package login

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var GitHubOauthConfig *oauth2.Config
const oauthStateGitHubString = "random" 

func InitGitHubOauthConfig() {
	GitHubOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	log.Println("✅ Loaded GitHub OAuth config")
	log.Println("CLIENT_ID =>", GitHubOauthConfig.ClientID)
	log.Println("REDIRECT_URI =>", GitHubOauthConfig.RedirectURL)
}

func GitHubLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := GitHubOauthConfig.AuthCodeURL(
		oauthStateGitHubString,
		oauth2.SetAuthURLParam("allow_signup", "true"),
		oauth2.SetAuthURLParam("login", "your_username"), 
	)
	
	log.Println("GitHubOauthConfig =>", GitHubOauthConfig)
	log.Println("RedirectURL =>", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("➡️ GitHubCallbackHandler started")

	if r.FormValue("state") != oauthStateGitHubString {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := GitHubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Code exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := GitHubOauthConfig.Client(context.Background(), token)

	// Get profile info
	userResp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()

	var user map[string]interface{}
	if err := json.NewDecoder(userResp.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode user", http.StatusInternalServerError)
		return
	}

	// Get email info
	emailResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		http.Error(w, "Failed to get email info", http.StatusInternalServerError)
		return
	}
	defer emailResp.Body.Close()

	var emails []map[string]interface{}
	if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
		http.Error(w, "Failed to decode emails", http.StatusInternalServerError)
		return
	}

	// Pick primary email
	primaryEmail := ""
	for _, email := range emails {
		if email["primary"].(bool) {
			primaryEmail = email["email"].(string)
			break
		}
	}

	user["email"] = primaryEmail
	user["verified_email"] = true

	// ⬇️ Respond with JSON like Google
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

