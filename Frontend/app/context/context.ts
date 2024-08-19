import React from "react";
import { TAsset, TUser } from "../utils/types";

export const dummyUser: TUser = {
  address: "",
  avatar: "",
  bio: "",
  createdAt: "",
  discord: "",
  email: "",
  id: "",
  isActive: false,
  isVerified: false,
  twitter: "",
  username: "",
  website: "",
};
export const dummyAsset: TAsset = {
  blockchain: "",
  category: "",
  createdAt: "",
  creatorId: "",
  description: "",
  id: "",
  issuanceDocumentsUrls: [],
  legalDocumentUrls: [],
  metadataUrl: "",
  name: "",
  properties: [],
  tokenId: "",
  totalSupply: 0,
  updatedAt: "",
  urls: [],
  symbol: "",
  favourites: 0,
  floorPrice: 0,
  isFavourite: false,
  isListedByUser: false,
  status: "unverified",
  tokenStandard: "",
  views: 0,
};

export const UserContext = React.createContext<{
  user: TUser;
  setUser: React.Dispatch<React.SetStateAction<TUser>>;
  signOut: () => void;
  address: string;
  copyAddress: () => void;
}>({
  user: dummyUser,
  setUser: () => {},
  signOut: () => {},
  address: "",
  copyAddress: () => {},
});

export const ThemeContext = React.createContext<{
  isDarkMode: boolean;
  setDarkMode: React.Dispatch<React.SetStateAction<boolean>>;
}>({
  isDarkMode: true,
  setDarkMode: () => {},
});

export const AssetContext = React.createContext<TAsset>(dummyAsset);
export const MediaQueryContext = React.createContext<boolean>(false);
