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
      keySeparator: ".",
      parseMissingKeyHandler: (key: string) => {
        const parts = key.split(".");

        // Handle special cases for objects and audio
        if (parts[0] === "object" || parts[0] === "audio") {
          return (
            parts[1]
              ?.split("_")
              .map(
                (word) =>
                  word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()
              )
              .join(" ") || key
          );
        }

        // For nested keys, try to make them more readable
        if (parts.length > 1) {
          const lastPart = parts[parts.length - 1];
          return lastPart
            .split("_")
            .map(
              (word) =>
                word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()
            )
            .join(" ");
        }

        // For single keys, just smart-capitalize and format
        return key
          .split("_")
          .map(
            (word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()
          )
          .join(" ");
      },
    });

  let node = document.createElement("div");
  document.body.appendChild(node);
  ReactDOM.createRoot(node).render(<App />);
};
