import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { api } from "./config";
import {
  TAsset,
  TAssetActivity,
  TCredentials,
  TNotifications,
  TNotificationPrefs,
  TOrderBooks,
  TOrderStatus,
  TOrderType,
  TResponse,
  TStatAssets,
} from "../types";
import { cookie } from "../cookie";

const NotificationsAPI = {
  markAllAsRead: ({
    onSuccess,
  }: {
    onSuccess: (res: TResponse<{ tokens: TCredentials }>) => void;
  }) =>
    useMutation({
      mutationKey: ["mark all as read"],
      mutationFn: () =>
        axios.put<{ tokens: TCredentials }>(
          api(`/notifications/read`),
          {},
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("credentials")["accessToken"]
              }`,
            },
          }
        ),
      onSuccess,
    }),

  getNotifications: () =>
    useMutation({
      mutationKey: ["get notifications"],
      mutationFn: ({ size = 1, page }: { size?: number; page: number }) =>
        axios<TNotifications>(api(`/notifications?page=${page}&size=${size}`), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
    }),
  getNotificationsPrefs: () =>
    useQuery({
      queryKey: ["get notifications prefs"],
      queryFn: () =>
        axios<TNotificationPrefs>(api(`/notifications/prefs`), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
    }),
  updateUserPrefs: () =>
    useMutation({
      mutationKey: ["update user prefs"],
      mutationFn: (data: any) =>
        axios.post<TNotificationPrefs>(api(`/notifications/prefs`), data, {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
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
        axios<TOrderBooks[]>(
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

export default NotificationsAPI;
