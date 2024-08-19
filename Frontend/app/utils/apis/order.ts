import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { api } from "./config";
import {
  TAsset,
  TCredentials,
  TErrorResponse,
  TResponse,
  TStatAssets,
  TUser,
} from "../types";
import { cookie } from "../cookie";

const OrderAPI = {
  create: ({
    onSuccess,
  }: {
    onSuccess: (res: TResponse<{ tokens: TCredentials }>) => void;
  }) =>
    useMutation({
      mutationKey: ["create order"],
      mutationFn: (data: {
        kind: string;
        quantity: number;
        assetId: string;
        type: "buy" | "sell";
        price?: number;
      }) =>
        axios<{ tokens: TCredentials }>(api(`/orders`), {
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
};

export default OrderAPI;
