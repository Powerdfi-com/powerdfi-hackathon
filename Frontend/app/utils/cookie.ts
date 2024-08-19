import Cookies from "js-cookie";
export const cookie = {
  set: (name: string, value: string) => Cookies.set(name, value),
  get: (name: string) => Cookies.get(name),
  setJson: (name: string, value: any) =>
    Cookies.set(name, JSON.stringify(value)),
  getJson: (name: string) =>
    Cookies.get(name) ? JSON.parse(Cookies.get(name)!) : null,
  remove: (name: string) => Cookies.remove(name),
};
