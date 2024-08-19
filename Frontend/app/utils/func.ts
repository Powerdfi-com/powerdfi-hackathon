import { cookie } from "./cookie";

export const months = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec",
];

export const generateArray = (l: number) => {
  return [...new Array(l)].map((_, index) => index);
};

export const dateToISO = (date: string) => {
  return new Date(date).toISOString();
};

export const formatDate = (date: Date) => {
  return (
    months[date.getMonth()] + " " + date.getDay() + ", " + date.getFullYear()
  );
};

export const isLinkActive = (url: string, path: string, isHome?: boolean) => {
  return url === path;
};

export const dateFromISO = (date: string) => new Date(date);

export const isDarkMode = cookie.get("theme")
  ? cookie.get("theme") === "dark"
  : true;
// window.matchMedia("(prefers-color-scheme: dark)").matches;

export const switchToDarkMode = () => cookie.set("theme", "dark");

export const switchToLightMode = () => cookie.set("theme", "light");
