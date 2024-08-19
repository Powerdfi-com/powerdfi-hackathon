import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { api } from "./config";
import {
  TAdmin,
  TCredentials,
  TErrorResponse,
  TResponse,
  TUser,
} from "../types";

const AuthAPI = {
  useFetchNonce: () =>
    useMutation({
      mutationKey: ["fetch nonce"],
      mutationFn: (address: string) =>
        axios<{ nonce: string }>(api(`/auth/${address}/nonce`)),
    }),
  verifySignature: () =>
    useMutation({
      mutationKey: ["sign message"],
      mutationFn: ({
        address,
        signature,
      }: {
        address: string;
        signature: string;
      }) =>
        axios<{ user: TUser; tokens: TCredentials }>(
          api(`/auth/${address}/verify`),
          {
            method: "POST",
            data: {
              signature,
            },
          }
        ),
    }),
  refreshToken: ({
    onSuccess,
    onError,
  }: {
    onSuccess: (res: TResponse<TCredentials>) => void;
    onError: TErrorResponse;
  }) =>
    useMutation({
      mutationKey: ["refresh token"],
      mutationFn: (refreshToken: string) =>
        axios<TCredentials>(api(`/auth/refresh-token`), {
          method: "POST",
          data: {
            refreshToken,
          },
        }),
      onSuccess,
      onError,
    }),
  adminSignIn: ({
    onSuccess,
  }: {
    onSuccess: (res: TResponse<{ data: TAdmin; tokens: TCredentials }>) => void;
  }) =>
    useMutation({
      mutationKey: ["sign in admin"],
      mutationFn: ({ email, password }: { email: string; password: string }) =>
        axios.post<{ data: TAdmin; tokens: TCredentials }>(api("/auth/admin"), {
          email,
          password,
        }),
      onSuccess,
    }),
};

export default AuthAPI;
