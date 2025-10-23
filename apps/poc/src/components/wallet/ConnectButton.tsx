import { useAccount, useConnect, useDisconnect } from "wagmi";
import { injected } from "wagmi/connectors";

export function WalletConnectButton() {
  const { isConnected, address } = useAccount();
  const { connect, connectors, status } = useConnect();
  const { disconnect } = useDisconnect();

  const metaMaskConnector =
    connectors.find((connector) => connector.id === "injected") ?? injected({ target: "metaMask" });

  if (isConnected) {
    return (
      <div className="flex items-center gap-3 rounded-full border border-emerald-400/30 bg-emerald-500/10 px-4 py-2">
        <span className="inline-flex h-2 w-2 rounded-full bg-emerald-400" />
        <span className="text-sm text-emerald-200">
          {address?.slice(0, 6)}...{address?.slice(-4)}
        </span>
        <button
          className="rounded-full border border-emerald-400/40 px-3 py-1 text-xs font-medium text-emerald-100 transition hover:bg-emerald-400/20"
          onClick={() => disconnect()}
        >
          Disconnect
        </button>
      </div>
    );
  }

  return (
    <button
      className="inline-flex items-center gap-2 rounded-full border border-slate-700 px-4 py-2 text-sm font-medium text-slate-100 transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60"
      disabled={status === "connecting"}
      onClick={() => connect({ connector: metaMaskConnector })}
    >
      <svg viewBox="0 0 24 24" className="h-4 w-4" aria-hidden>
        <path
          fill="currentColor"
          d="M20.1 3.8 13.7 8l1.2-2.8 5.2-1.4Zm-.2 14 .2-2.4-3.7.8 3.5 1.6Zm-15.8 0 3.5-1.6-3.7-.8.2 2.4Zm.2-14 5.2 1.4L10.7 8 4.3 3.8ZM12 9.5l-2.2 4.6 2.2 1.6 2.2-1.6L12 9.5Zm8.7 3.8 1-3.2-2.5-.2 1.5 3.4Zm-17.4 0L4.8 10l-2.5.2 1.5 3.4Zm7.4 6.3-3.8-1.9L9 15.8l1.7 3.8Zm2.6 0 1.7-3.8 2.1 1.9-3.8 1.9Z"
        />
      </svg>
      {status === "connecting" ? "Connecting..." : "Connect MetaMask"}
    </button>
  );
}
