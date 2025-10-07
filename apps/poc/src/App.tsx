import { http, createConfig, WagmiProvider } from "wagmi";
import { sepolia } from "wagmi/chains";
import { injected } from "wagmi/connectors";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { EventDashboard } from "./components/EventDashboard";
import { WalletProvider } from "./context/WalletContext";
import { Layout } from "./components/Layout";

const queryClient = new QueryClient();

const wagmiConfig = createConfig({
  chains: [sepolia],
  transports: {
    [sepolia.id]: http(),
  },
  connectors: [injected({ target: "metaMask" })],
});

export function App() {
  return (
    <WagmiProvider config={wagmiConfig}>
      <QueryClientProvider client={queryClient}>
        <WalletProvider>
          <Layout>
            <EventDashboard />
          </Layout>
        </WalletProvider>
      </QueryClientProvider>
    </WagmiProvider>
  );
}
