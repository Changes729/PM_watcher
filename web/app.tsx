import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route, useNavigate } from "react-router";
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import HttpBackend from "i18next-http-backend";

import Homepage from "./pages/home";

import "./app.scss";

function App() {
  onload = () => {
    useNavigate()(document.location.pathname);
  };

  return (
    <BrowserRouter>
      <Routes>
        <Route index path="/" element={<Homepage />} />
      </Routes>
    </BrowserRouter>
  );
}

/** Main Start */
window.onload = () => {
  i18n
    .use(initReactI18next)
    .use(HttpBackend)
    .init({
      fallbackLng: "zh",

      backend: {
        loadPath: "/locales/{{lng}}/{{ns}}.json",
      },

      ns: ["common"],
      defaultNS: "common",
    });

  let node = document.createElement("div");
  document.body.appendChild(node);
  ReactDOM.createRoot(node).render(<App />);
};
