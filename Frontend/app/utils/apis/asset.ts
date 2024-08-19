import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { api } from "./config";
import {
  TAsset,
  TAssetActivity,
  TCredentials,
  TOrderBooks,
  TOrderStatus,
  TOrderType,
  TResponse,
  TStatAssets,
} from "../types";
import { cookie } from "../cookie";

const AssetAPI = {
  create: ({
    onSuccess,
  }: {
    onSuccess: (res: TResponse<{ tokens: TCredentials }>) => void;
  }) =>
    useMutation({
      mutationKey: ["create asset"],
      mutationFn: (data: {
        name: string;
        symbol: string;
        totalSupply: number;
        description: string;
        urls: string[];
        legalDocs: string[];
        issuanceDocs: string[];
        categoryId: number;
        properties: string;
        blockchainId: string;
      }) =>
        axios<{ tokens: TCredentials }>(api(`/assets`), {
          method: "POST",
          data,
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
      onSuccess,
    }),

  getAssetById: (id: string) =>
    useQuery({
      queryKey: ["get asset by id", id],
      queryFn: () =>
        axios<TAsset>(api(`/assets/${id}`), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
    }),
  getRecommendedAssets: (id: string) =>
    useQuery({
      queryKey: ["get recommended assets", id],
      queryFn: () =>
        axios<TStatAssets>(api(`/assets/${id}/recommended`), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
    }),
  getCategories: () =>
    useQuery({
      queryKey: ["get categories"],
      queryFn: () =>
        axios<{ id: number; name: string; slug: string }[]>(
          api(`/assets/categories`),
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("credentials")["accessToken"]
              }`,
            },
          }
        ),
    }),
  getChains: () =>
    useQuery({
      queryKey: ["get chains"],
      queryFn: () =>
        axios<{ id: number; name: string; logo: string }[]>(
          api(`/assets/chains`),
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("credentials")["accessToken"]
              }`,
            },
          }
        ),
    }),
  getActivities: ({ assetId }: { assetId: string }) =>
    useQuery({
      queryKey: ["get activities"],
      queryFn: () =>
        axios<TAssetActivity[]>(
          api(`/assets/${assetId}/activities?page=1&size=9`),
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("credentials")["accessToken"]
              }`,
            },
          }
        ),
    }),
  getOrderBook: ({
    id,
    status,
    type,
  }: {
    id: string;
    status?: TOrderStatus;
    type?: TOrderType;
  }) =>
    useQuery({
      queryKey: ["get order book", { id, status, type }],
      queryFn: () =>
        axios<TOrderBooks>(
          api(
            `/assets/${id}/orderbook?${status && `status=${status}`}${
              type && `&type=${type}`
            }`
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

export default AssetAPI;
