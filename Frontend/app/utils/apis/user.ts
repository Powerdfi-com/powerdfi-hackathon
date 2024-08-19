import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { api } from "./config";
import {
  TAsset,
  TAssetActivity,
  TCredentials,
  TErrorResponse,
  TResponse,
  TStatAssets,
  TUser,
} from "../types";
import { cookie } from "../cookie";

const UserAPI = {
  activateUser: ({
    onSuccess,
  }: {
    onSuccess: (res: TResponse<{ tokens: TCredentials }>) => void;
  }) =>
    useMutation({
      mutationKey: ["activate user"],
      mutationFn: (data: { username: string; email: string }) =>
        axios<{ tokens: TCredentials }>(api(`/user/activate`), {
          method: "PATCH",
          data,
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
      onSuccess,
    }),
  getProfile: () =>
    useMutation({
      mutationKey: ["fetch profile"],
      mutationFn: () =>
        axios<TUser>(api("/user/me"), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
    }),
  getKycLink: ({
    onSuccess,
  }: {
    onSuccess: (res: TResponse<{ link: string }>) => void;
  }) =>
    useMutation({
      mutationKey: ["get kyc link"],
      mutationFn: () =>
        axios<{ link: string }>(api(`/user/kyc-link`), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
      onSuccess,
    }),
  updateProfile: ({
    onSuccess,
  }: {
    onSuccess: (res: TResponse<TUser>) => void;
  }) =>
    useMutation({
      mutationKey: ["update profile"],
      mutationFn: (data: {
        bio: string;
        website: string;
        discord: string;
        twitter: string;
        avatar: string;
      }) =>
        axios<TUser>(api(`/user/me`), {
          method: "PATCH",
          data,
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
      onSuccess,
    }),

  getCreatedAssets: (id: string) =>
    useQuery({
      queryKey: ["get created assets"],
      queryFn: () => axios<TAsset[]>(api(`/user/${id}/created`)),
    }),
  getListedAssets: (id: string) =>
    useQuery({
      queryKey: ["get listed assets"],
      queryFn: () => axios<TAsset[]>(api(`/user/${id}/listings`)),
    }),
  getWalletDetails: () =>
    useQuery({
      queryKey: ["get wallet details"],
      queryFn: () =>
        axios<{ accoundId: string; address: string; balance: number }>(
          api(`/user/wallet`),
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("credentials")["accessToken"]
              }`,
            },
          }
        ),
    }),
  getPortfolio: () =>
    useQuery({
      queryKey: ["get portfolio"],
      queryFn: () =>
        axios<TAsset[]>(api("/user/portfolio"), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
    }),
  getActivities: () =>
    useQuery({
      queryKey: ["get portfolio"],
      queryFn: () =>
        axios<{ activities: TAssetActivity[]; total: number }>(
          api("/user/activities?page=1&size=5"),
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("credentials")["accessToken"]
              }`,
            },
          }
        ),
    }),
  getTopAssets: () =>
    useMutation({
      mutationKey: ["get user top assets"],
      mutationFn: ({
        page,
        size,
        range,
        categoryId,
        blockchain,
      }: {
        page?: number;
        size?: number;
        range?: number;
        categoryId?: number;
        blockchain?: string;
      }) =>
        axios<TStatAssets>(
          api(
            `/user/top-assets/perfs?${page ? `page=${page}&` : ""}${
              size ? `size=${size}&` : ""
            }${range ? `range=${range}&` : ""}${
              categoryId ? `categoryId=${categoryId}&` : ""
            }${blockchain ? `blockchain=${blockchain}` : ""}`
          ),
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("credentials")["accessToken"]
              }`,
            },
          }
        ),
    }),
  getTrendingAssets: () =>
    useMutation({
      mutationKey: ["get user trending assets"],
      mutationFn: ({
        page,
        size,
        range,
        categoryId,
        blockchain,
      }: {
        page?: number;
        size?: number;
        range?: number;
        categoryId?: number;
        blockchain?: string;
      }) =>
        axios<TStatAssets>(
          api(
            `/user/trending-assets/perfs?${page ? `page=${page}&` : ""}${
              size ? `size=${size}&` : ""
            }${range ? `range=${range}&` : ""}${
              categoryId ? `categoryId=${categoryId}&` : ""
            }${blockchain ? `blockchain=${blockchain}` : ""}`
          ),
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("credentials")["accessToken"]
              }`,
            },
          }
        ),
    }),
};

export default UserAPI;
