import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { api } from "./config";
import { TAdmin, TAsset, TResponse, TStatAssets } from "../types";
import { cookie } from "../cookie";

const AdminAPI = {
  getAdmins: () =>
    useMutation({
      mutationKey: ["fetch admins"],
      mutationFn: (page: number) =>
        axios<{ total: number; admins: TAdmin[] }>(
          api(`/admin?page=${page}&size=2`),
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("adminCredentials")["accessToken"]
              }`,
            },
          }
        ),
    }),
  getCreatorsSurvey: () =>
    useQuery({
      queryKey: ["fetch creators survey"],
      queryFn: () =>
        axios<
          { month: number; heritageUsersCount: number; newUsersCount: number }[]
        >(api(`/admin/survey/creators`), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("adminCredentials")["accessToken"]
            }`,
          },
        }),
    }),
  getAdminById: (id: string) =>
    useQuery({
      queryKey: ["fetch admin", { id }],
      queryFn: () =>
        axios<TAdmin>(api(`/admin/${id}`), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("adminCredentials")["accessToken"]
            }`,
          },
        }),
    }),
  getCategoriesSurvey: () =>
    useQuery({
      queryKey: ["fetch categories survey"],
      queryFn: () =>
        axios<any[]>(api(`/admin/survey/categories`), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("adminCredentials")["accessToken"]
            }`,
          },
        }),
    }),
  getAssetsSurvey: () =>
    useQuery({
      queryKey: ["fetch assets survey"],
      queryFn: () =>
        axios<
          {
            month: number;
            verifiedAssetsCount: number;
            unverifiedAssetsCount: number;
          }[]
        >(api(`/admin/survey/verified-assets`), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("adminCredentials")["accessToken"]
            }`,
          },
        }),
    }),
  getAssets: () =>
    useQuery({
      queryKey: ["fetch assets"],
      queryFn: () =>
        axios<TAsset[]>(api("/admin/assets"), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("adminCredentials")["accessToken"]
            }`,
          },
        }),
    }),
  getStats: () =>
    useQuery({
      queryKey: ["fetch stats"],
      queryFn: () =>
        axios<{
          usersCount: number;
          creatorsCount: number;
          assetsCount: number;
          percentageChangeCreators: string;
          percentageChangeUsers: string;
        }>(api("/admin/stats"), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("adminCredentials")["accessToken"]
            }`,
          },
        }),
    }),
  getActivities: () =>
    useMutation({
      mutationKey: ["fetch activities"],
      mutationFn: (page: number) =>
        axios(api("/admin/stats"), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("adminCredentials")["accessToken"]
            }`,
          },
        }),
    }),
  getRoles: () =>
    useQuery({
      queryKey: ["fetch assets"],
      queryFn: () =>
        axios<{ id: number; name: string }[]>(api("/admin/roles"), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("adminCredentials")["accessToken"]
            }`,
          },
        }),
    }),

  getNotificationsPrefs: () =>
    useQuery({
      queryKey: ["fetch notfications prefs"],
      queryFn: () =>
        axios<{ login: boolean; created: boolean }>(
          api("/admin/notifications/prefs"),
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("adminCredentials")["accessToken"]
              }`,
            },
          }
        ),
    }),

  updateAdminPrefs: () =>
    useMutation({
      mutationKey: ["update admin prefs"],
      mutationFn: (data: any) =>
        axios.post<{ created: boolean; login: boolean }>(
          api(`/admin/notifications/prefs`),
          data,
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("adminCredentials")["accessToken"]
              }`,
            },
          }
        ),
    }),
  approveAsset: ({ onSuccess }: { onSuccess: () => void }) =>
    useMutation({
      mutationKey: ["verify asset"],
      mutationFn: (id: string) =>
        axios.post(
          api("/admin/approve"),
          {
            assetId: id,
          },
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("adminCredentials")["accessToken"]
              }`,
            },
          }
        ),
      onSuccess,
    }),
  rejectAsset: ({ onSuccess }: { onSuccess: () => void }) =>
    useMutation({
      mutationKey: ["reject asset"],
      mutationFn: (id: string) =>
        axios.post(
          api("/admin/reject"),
          {
            assetId: id,
          },
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("adminCredentials")["accessToken"]
              }`,
            },
          }
        ),
      onSuccess,
    }),
  createAdmin: () =>
    useMutation({
      mutationKey: ["admin create admin"],
      mutationFn: ({
        email,
        password,
        role,
      }: {
        email: string;
        password: string;
        role: number;
      }) =>
        axios.post(
          api("/admin"),
          {
            email,
            password,
            role,
          },
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("adminCredentials")["accessToken"]
              }`,
            },
          }
        ),
    }),
  updateAdmin: () =>
    useMutation({
      mutationKey: ["admin update admin"],
      mutationFn: ({
        id,
        email,
        password,
        role,
      }: {
        id: string;
        email: string;
        password: string;
        role: number;
      }) =>
        axios.patch(
          api(`/admin/${id}`),
          {
            email,
            password,
            role,
          },
          {
            headers: {
              Authorization: `Bearer ${
                cookie.getJson("adminCredentials")["accessToken"]
              }`,
            },
          }
        ),
    }),
};

export default AdminAPI;
