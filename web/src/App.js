import { Router } from "@reach/router";

import Main from "Main";
import Home from "pages/Home";
import Login from "pages/Login";

import React, { useCallback, useState } from "react";
import { ProtectedRoute, NotFound } from "shared/components";
import { AuthContext } from "shared/contexts";

import "./App.less";

const LS_USER_KEY = "user";

function getUserFromStorage() {
  let user;
  try {
    user = JSON.parse(localStorage.getItem(LS_USER_KEY)) || undefined;
  } catch (err) {
    console.err(err);
  }
  return user;
}

function App() {
  const [user, setUser] = useState(getUserFromStorage);

  const login = useCallback((user) => {
    localStorage.setItem(LS_USER_KEY, JSON.stringify(user));
    setUser(user);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem(LS_USER_KEY);
    setUser(undefined);
  }, []);

  return (
    <AuthContext.Provider value={[user, login, logout]}>
      <Router id="router">
        <Login path="login" />
        <ProtectedRoute user={user} component={Main} path="/">
          <Home path="/" />
          <NotFound default />
        </ProtectedRoute>
        <NotFound default />
      </Router>
    </AuthContext.Provider>
  );
}

export default App;
