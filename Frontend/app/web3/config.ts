import { defaultWagmiConfig } from "@web3modal/wagmi";
import { hederaTestnet } from "wagmi/chains";
import { createWeb3Modal } from '@web3modal/wagmi/react'

const projectId = "a3c24e6e2c5988b2ccc07499855f681a";

export const config = defaultWagmiConfig({
  chains: [hederaTestnet],
  projectId,
  metadata: {
    name: "Web3Modal",
    description: "Web3Modal Example",
    url: "https://web3modal.com", // origin must match your domain & subdomain
    icons: ["https://avatars.githubusercontent.com/u/37784886"],
  },
});


createWeb3Modal({
  wagmiConfig: config,
  projectId,
  enableAnalytics: true, // Optional - defaults to your Cloud configuration
  enableOnramp: true // Optional - false as default
})
