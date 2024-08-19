import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { api } from "./config";
import { TStatAssets } from "../types";
import { cookie } from "../cookie";

const StatsAPI = {
  getTopAssets: () =>
    useMutation({
      mutationKey: ["get top assets"],
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
            `/stats/top-assets?${page ? `page=${page}&` : ""}${
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
      mutationKey: ["get trending assets"],
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
            `/stats/trending-assets?${page ? `page=${page}&` : ""}${
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
  getTopAssetsPerformance: ({
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
    useMutation({
      mutationKey: [
        "get top assets performance",
        page,
        size,
        range,
        categoryId,
        blockchain,
      ],
      mutationFn: () =>
        axios<TStatAssets>(
          api(
            `/stats/top-assets/perfs?${page ? `page=${page}&` : ""}${
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
  getTrendingAssetsPerformance: () =>
    useMutation({
      mutationKey: ["get trending assets performance"],
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
            `/stats/trending-assets/perfs?${page ? `page=${page}&` : ""}${
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

export default StatsAPI;
