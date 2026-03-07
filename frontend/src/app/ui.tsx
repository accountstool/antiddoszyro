import { createContext, useContext, useEffect, useMemo, useState, type PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";

type Theme = "light" | "dark";

type UIContextValue = {
  theme: Theme;
  toggleTheme: () => void;
  language: string;
  setLanguage: (language: string) => void;
};

const UIContext = createContext<UIContextValue | null>(null);

export function UIProvider({ children }: PropsWithChildren) {
  const { i18n } = useTranslation();
  const [theme, setTheme] = useState<Theme>(() => {
    const saved = localStorage.getItem("shieldpanel-theme");
    return saved === "dark" ? "dark" : "light";
  });
  const [language, setLanguageState] = useState(() => localStorage.getItem("shieldpanel-language") || "en");

  useEffect(() => {
    document.documentElement.classList.toggle("theme-dark", theme === "dark");
    localStorage.setItem("shieldpanel-theme", theme);
  }, [theme]);

  useEffect(() => {
    void i18n.changeLanguage(language);
    localStorage.setItem("shieldpanel-language", language);
  }, [i18n, language]);

  const value = useMemo<UIContextValue>(
    () => ({
      theme,
      language,
      toggleTheme() {
        setTheme((current) => (current === "dark" ? "light" : "dark"));
      },
      setLanguage(nextLanguage: string) {
        setLanguageState(nextLanguage);
      }
    }),
    [language, theme]
  );

  return <UIContext.Provider value={value}>{children}</UIContext.Provider>;
}

export function useUI() {
  const context = useContext(UIContext);
  if (!context) {
    throw new Error("useUI must be used within UIProvider");
  }
  return context;
}
