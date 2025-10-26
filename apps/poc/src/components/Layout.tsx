import { ReactNode } from "react";
import { WalletConnectButton } from "./wallet/ConnectButton";

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">
      <div className="pointer-events-none fixed inset-0 -z-20 overflow-hidden">
        <div className="absolute left-1/2 top-0 h-[480px] w-[480px] -translate-x-1/2 rounded-full bg-emerald-500/10 blur-[160px]" />
        <div className="absolute bottom-0 right-0 h-[380px] w-[420px] rounded-full bg-indigo-500/10 blur-[140px]" />
      </div>

      <header className="border-b border-slate-800/70 bg-slate-950/70 backdrop-blur">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-5">
          <div>
            <div className="flex items-center gap-2 text-sm text-emerald-300">
              <span className="inline-flex h-2 w-2 rounded-full bg-emerald-400" />
              Sepolia Network
            </div>
            <h1 className="text-2xl font-semibold tracking-tight">DECM Event Portal</h1>
            <p className="text-sm text-slate-400">
              Manage hosts, participants, tickets, and certificates in one dashboard
            </p>
          </div>
          <WalletConnectButton />
        </div>
      </header>

      <main className="mx-auto flex max-w-6xl flex-col gap-8 px-6 py-10">{children}</main>
    </div>
  );
}
