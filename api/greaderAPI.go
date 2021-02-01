package api

import (
  "net/http"
  log "github.com/sirupsen/logrus"
  "strings"
  "encoding/json"
  "github.com/gorilla/mux"
)

type GreaderAPIAccountsClientLoginResponse struct {
  SID                 string          `json:"SID,omitempty"`
  LSID                string          `json:"LSID,omitempty"`
  Auth                string          `json:"Auth,omitempty"`
}

type GreaderAPIUserInfoResponse struct {
  UserID              string          `json:"userId,omitempty"`
  UserName            string          `json:"userName,omitempty"`
  UserProfileId       string          `json:"userProfileId,omitempty"`
  UserEmail           string          `json:"userEmail,omitempty"`
  IsBloggerUser       bool            `json:"isBloggerUser,omitempty"`
  SignupTimeSec       int             `json:"signupTimeSec,omitempty"`
  IsMultiLoginEnabled bool            `json:"isMultiLoginEnabled,omitempty"`
  IsPremium           bool            `json:"isPremium,omitempty"`
}

func greaderAPIHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  if r.Method == http.MethodOptions {
      return
  }

  // TODO: Implement Google Reader API
  // See:
  // - http://code.google.com/p/pyrfeed/wiki/GoogleReaderAPI
  // - https://blog.martindoms.com/2009/10/16/using-the-google-reader-api-part-2
  // - https://ranchero.com/downloads/GoogleReaderAPI-2009.pdf
  // - https://github.com/theoldreader/api
  // - https://github.com/devongovett/reader
  // - https://github.com/FreshRSS/FreshRSS/blob/master/p/api/greader.php

  w.WriteHeader(http.StatusNoContent)
  return
}

func greaderAPIAccountsClientLoginHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  if r.Method == http.MethodOptions {
      return
  }

  r.ParseMultipartForm(8192)

  email := r.FormValue("Email")
  if email == "" {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  passwd := r.FormValue("Passwd")
  if passwd == "" {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }


  accountsClientLoginResponse := GreaderAPIAccountsClientLoginResponse{
    SID: "none",
    LSID: "none",
    Auth: GetApiKey(email, passwd),
  }

  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(accountsClientLoginResponse)
  return
}

func greaderAPIUserInfo(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  if r.Method == http.MethodOptions {
      return
  }

  authHeader := r.Header.Get("Authorization")
  authHeaderSplit := strings.Split(authHeader, " ")
  authSplit := strings.Split(authHeaderSplit[1], "=")
  authHash := authSplit[1]
  log.Debug(authHash)

  userInfoResponse := GreaderAPIUserInfoResponse{
    UserID: "123",
    UserName: "Test",
    UserProfileId: "123",
    UserEmail: "test@test.com",
    IsBloggerUser: false,
    SignupTimeSec: 0,
    IsMultiLoginEnabled: false,
    IsPremium: true,
  }

  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(userInfoResponse)
  return
}

func greaderAPI(r *mux.Router) {
  r.HandleFunc("/", greaderAPIHandler)
  r.HandleFunc("/accounts/ClientLogin", greaderAPIAccountsClientLoginHandler)
  r.HandleFunc("/reader/api/0/user-info", greaderAPIUserInfo)
}
