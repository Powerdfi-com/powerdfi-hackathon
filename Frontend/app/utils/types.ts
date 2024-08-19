import { AxiosError } from "axios";

type TProvider = {
  name: string;
  image: string;
};

export interface TResponse<T> {
  data: T;
}

export type TErrorResponse = (
  res: AxiosError<{ responseMessage: string }>
) => void;

export type TFilter = {
  type: string;
  blockchain: string;
  range: number;
  categoryId?: number;
};

export type TUser = {
  id: string;
  address: string;
  isActive: false;
  isVerified: false;
  avatar: string;
  bio: string;
  twitter: string;
  discord: string;
  website: string;
  createdAt: string;
  username: string;
  email: string;
};

export type TCredentials = {
  accessToken: string;
  refreshToken: string;
  expiresAt: number;
};

export type TAdmin = {
  id: string;
  email: string;
  roles: string[];
  createdAt: string;
};

export type TAsset = {
  id: string;
  name: string;
  description: string;
  category: string;
  creatorId: string;
  blockchain: string;
  tokenId: string;
  metadataUrl: string;
  totalSupply: number;
  properties: string[];
  urls: string[];
  legalDocumentUrls: string[];
  issuanceDocumentsUrls: string[];
  createdAt: string;
  updatedAt: string;
  symbol: string;
  favourites: number;
  views: number;
  isFavourite: boolean;
  status: TAssetStatus;
  isListedByUser: boolean;
  floorPrice: number;
  tokenStandard: string;
};

export type TAssetActivity = {
  id: string;
  action: "sale" | "mint";
  assetId: string;
  assetName: string;
  price: number;
  currency: string;
  quantity: number;
  fromUserId: string;
  toUserId: string;
  createdAt: string;
};

export type TAssetStatus = "verified" | "unverified" | "rejected";

export type TNotifications = {
  notifications: {
    id: string;
    type: "reject" | "approve" | "sale";
    userId: string;
    data: {
      assetId: string;
      assetName: string;
      reason: string;
    };
    createdAt: string;
    viewed: boolean;
  }[];
  total: number;
};

export type TNotificationPrefs = {
  sale: boolean;
  verified: boolean;
  rejected: boolean;
  login: boolean;
};

export type TStatAssets = {
  assets: {
    assetId: string;
    category: string;
    categoryName: string;
    name: string;
    blockchain: string;
    logo: string;
    volume: number;
    owners: number;
    floorPrice: number;
    blockchainLogo: string;
    creatorUsername: string;
    status: TAssetStatus;
    isVerified: boolean;
    percentageChange: string;
    priceChanges?: {
      timestamp: string;
      price: number;
    }[];
  }[];
  total: number;
};

export type TOrderBooks = {
  orders: {
    id: string;
    assetId: string;
    userId: string;
    type: TOrderType;
    kind: "market";
    price: number;
    quantity: number;
    status: TOrderStatus;
    createdAt: string;
    updatedAt: string;
  }[];
  total: number;
};

export type TOrderType = "buy" | "sell";
export type TOrderStatus = "open" | "cancelled" | "filled" | "partial";
