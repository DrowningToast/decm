import { useState } from "react";
import type { Address, Hex } from "viem";
import { isHex } from "viem";
import { useAccount, usePublicClient, useWalletClient } from "wagmi";
import { sepolia } from "wagmi/chains";
import { EventABI } from "../abi/Event";
import { EventTicketABI } from "../abi/EventTicket";
import { EventCertificateABI } from "../abi/EventCertificate";
import {
  chainId,
  EVENT_ADDRESS,
  EVENT_CERTIFICATE_ADDRESS,
  EVENT_TICKET_ADDRESS,
} from "../lib/addresses";
import { useWalletContext } from "../context/WalletContext";
import { useParticipantMessage } from "../hooks/useParticipants";

interface EventParticipantPanelProps {
  eventAddress: Address;
}

export function EventParticipantPanel({
  eventAddress = EVENT_ADDRESS,
}: EventParticipantPanelProps) {
  const { address } = useAccount();
  const { refreshRoles } = useWalletContext();
  const publicClient = usePublicClient({ chainId });
  const { data: walletClient } = useWalletClient({ chainId: sepolia.id });
  const [signature, setSignature] = useState<Hex>("0x");
  const [tokenId, setTokenId] = useState<number>(0);
  const messageQuery = useParticipantMessage((address ?? "0x") as Address, eventAddress);

  const write = async (config: {
    address: Address;
    abi: unknown;
    functionName: string;
    args: unknown[];
  }) => {
    if (!walletClient || !address) throw new Error("Wallet client unavailable");
    const hash = await walletClient.writeContract({
      account: address,
      ...config,
    });
    await publicClient?.waitForTransactionReceipt({ hash });
  };

  const handleLeaveEvent = async () => {
    if (!address) return;
    await write({
      address: eventAddress,
      abi: EventABI,
      functionName: "leaveEvent",
      args: [address],
    });
    await refreshRoles();
  };

  const handleSignMessage = async () => {
    if (!walletClient || !address) return;
    const message = messageQuery.data;
    if (!message || message === "0x") {
      alert("No signing message has been set yet.");
      return;
    }
    const signed = await walletClient.signMessage({
      account: address,
      message: { raw: message },
    });
    setSignature(signed as Hex);
  };

  const handleVerifySignature = async () => {
    if (!address) return;
    if (!signature || !isHex(signature)) {
      alert("Provide a valid signature");
      return;
    }
    const recovered = await publicClient?.readContract({
      address: eventAddress,
      abi: EventABI,
      functionName: "verifySignature",
      args: [address, signature],
    });
    alert(`Recovered address: ${recovered}`);
  };

  const handleViewTicket = async () => {
    if (!EVENT_TICKET_ADDRESS) {
      alert("Ticket contract address not configured.");
      return;
    }
    const uri = await publicClient?.readContract({
      address: EVENT_TICKET_ADDRESS,
      abi: EventTicketABI,
      functionName: "tokenURI",
      args: [BigInt(tokenId)],
    });
    alert(uri);
  };

  const handleViewCertificate = async () => {
    if (!EVENT_CERTIFICATE_ADDRESS) {
      alert("Certificate contract address not configured.");
      return;
    }
    const uri = await publicClient?.readContract({
      address: EVENT_CERTIFICATE_ADDRESS,
      abi: EventCertificateABI,
      functionName: "tokenURI",
      args: [BigInt(tokenId)],
    });
    alert(uri);
  };

  return (
    <div className="overflow-hidden rounded-3xl border border-slate-800/60 bg-slate-950/70 shadow-2xl shadow-slate-950/40">
      <div className="relative">
        <div className="absolute inset-0 bg-gradient-to-br from-slate-950 via-indigo-950/20 to-slate-950" />
        <div className="relative grid gap-8 px-8 py-8">
          <div className="flex flex-col gap-3">
            <span className="inline-flex w-fit items-center gap-2 rounded-full border border-indigo-500/30 bg-indigo-500/15 px-3 py-1 text-xs text-indigo-200">
              <span className="inline-flex h-1.5 w-1.5 rounded-full bg-indigo-300" /> Participant
              Panel
            </span>
            <h2 className="text-xl font-semibold tracking-tight text-slate-100">
              Participant Tools
            </h2>
            <p className="text-sm text-slate-400">
              Sign your attendance, verify on-chain signatures, and retrieve your ticket or
              certificate metadata directly from Sepolia.
            </p>
          </div>

          <section className="grid gap-4 rounded-2xl border border-slate-800/60 bg-slate-950/70 p-5">
            <div className="flex flex-col gap-2 text-sm">
              <span className="text-xs uppercase tracking-widest text-slate-500">
                Signing Message
              </span>
              <code className="rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-xs text-slate-300">
                {messageQuery.data ?? "No message available"}
              </code>
            </div>
            <div className="flex flex-wrap items-center gap-3">
              <button
                className="rounded-lg bg-emerald-500/15 px-4 py-2 text-sm font-medium text-emerald-200 transition hover:bg-emerald-500/25"
                onClick={handleSignMessage}
              >
                Sign Message
              </button>
              <input
                className="w-full max-w-md rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-xs text-slate-100 focus:border-sky-400/50 focus:outline-none focus:ring-2 focus:ring-sky-400/20"
                value={signature}
                onChange={(event) => setSignature(event.target.value as Hex)}
                placeholder="0x signature"
              />
              <button
                className="rounded-lg bg-sky-500/15 px-4 py-2 text-sm font-medium text-sky-200 transition hover:bg-sky-500/25"
                onClick={handleVerifySignature}
              >
                Verify Signature
              </button>
            </div>
          </section>

          <section className="grid gap-4 rounded-2xl border border-slate-800/60 bg-slate-950/70 p-5">
            <h3 className="text-sm font-medium text-slate-200">Attendance</h3>
            <button
              className="rounded-lg bg-rose-500/15 px-4 py-2 text-sm font-medium text-rose-200 transition hover:bg-rose-500/25"
              onClick={handleLeaveEvent}
            >
              Leave Event
            </button>
          </section>

          <section className="grid gap-4 rounded-2xl border border-slate-800/60 bg-slate-950/70 p-5">
            <h3 className="text-sm font-medium text-slate-200">Credential Metadata</h3>
            <div className="grid gap-2">
              <label className="text-xs uppercase tracking-widest text-slate-500">Token ID</label>
              <input
                className="w-full max-w-xs rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100 focus:border-indigo-400/50 focus:outline-none focus:ring-2 focus:ring-indigo-400/20"
                type="number"
                value={tokenId}
                onChange={(event) => setTokenId(Number(event.target.value))}
              />
            </div>
            <div className="flex flex-wrap gap-3">
              <button
                className="rounded-lg bg-indigo-500/15 px-4 py-2 text-sm font-medium text-indigo-200 transition hover:bg-indigo-500/25"
                onClick={handleViewTicket}
              >
                View Ticket Metadata
              </button>
              <button
                className="rounded-lg bg-purple-500/15 px-4 py-2 text-sm font-medium text-purple-200 transition hover:bg-purple-500/25"
                onClick={handleViewCertificate}
              >
                View Certificate Metadata
              </button>
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}
