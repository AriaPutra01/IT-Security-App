import React, { createContext, useContext, useEffect, useState } from "react";
import { useCookies } from "react-cookie";
import { jwtDecode } from "jwt-decode";

const TokenContext = createContext();

export const TokenProvider = ({ children }) => {
  const [cookies] = useCookies(["token"]);
  const [token, setToken] = useState(cookies.token || null);
  const [userDetails, setUserDetails] = useState({});

  useEffect(() => {
    if (token) {
      const decoded = jwtDecode(token);
      setUserDetails({
        username: decoded.username,
        email: decoded.email,
        role: decoded.role,
      });
    } else {
      setUserDetails({});
    }
  }, [token]);

  return (
    <TokenContext.Provider value={{ token, userDetails, setToken }}>
      {children}
    </TokenContext.Provider>
  );
};

export const useToken = () => {
  return useContext(TokenContext);
};
