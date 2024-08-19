import { useMutation, useQuery } from "@tanstack/react-query";
import axios from "axios";
import { api } from "./config";
import { TStatAssets } from "../types";
import { cookie } from "../cookie";

const ListingAPI = {
  createListing: ({ onSuccess }: { onSuccess: () => void }) =>
    useMutation({
      mutationKey: ["create listing"],
      mutationFn: ({
        assetId,
        currency,
        price,
        quantity,
        startAt,
        endAt,
        min_investment_amount,
        max_investment_amount,
        max_raise_amount,
        min_raise_amount,
        type,
      }: {
        assetId: string;
        currency: string[];
        price: number;
        quantity: number;
        startAt: string;
        endAt: string;
        min_investment_amount: number;
        max_investment_amount: number;
        max_raise_amount: number;
        min_raise_amount: number;
        type: string;
      }) =>
        axios<TStatAssets[]>(api(`/listings`), {
          method: "POST",
          data: {
            assetId,
            currency,
            price,
            quantity,
            startAt,
            endAt,
            min_investment_amount,
            max_investment_amount,
            max_raise_amount,
            min_raise_amount,
            type,
          },
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
      onSuccess,
    }),
  getChainTokens: (chainId: string) =>
    useQuery({
      queryKey: ["fetch chain tokens"],
      queryFn: () =>
        axios<
          {
            id: string;
            name: string;
          }[]
        >(api(`/listings/${chainId}/tokens`), {
          headers: {
            Authorization: `Bearer ${
              cookie.getJson("credentials")["accessToken"]
            }`,
          },
        }),
    }),
};

export default ListingAPI;
