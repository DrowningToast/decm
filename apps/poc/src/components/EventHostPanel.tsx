import { useState } from "react";
import { type Address, type Hex, isHex } from "viem";
import { useAccount, usePublicClient, useWalletClient } from "wagmi";
import { sepolia } from "wagmi/chains";
import { EventABI } from "../abi/Event";
import { EventAccessManagerABI } from "../abi/EventAccessManager";
import { EventTicketABI } from "../abi/EventTicket";
import { EventCertificateABI } from "../abi/EventCertificate";
import { EventStatus } from "../types/Event";
import {
  ACCESS_MANAGER_ADDRESS,
  EVENT_ADDRESS,
  EVENT_TICKET_ADDRESS,
  EVENT_CERTIFICATE_ADDRESS,
  chainId,
} from "../lib/addresses";
import { useEventData } from "../hooks/useEventData";
import { useParticipantMessage } from "../hooks/useParticipants";
import { useWalletContext } from "../context/WalletContext";

interface EventHostPanelProps {
  eventAddress?: Address;
  accessManagerAddress?: Address;
}

const defaultParticipant = "0x" as Address;

export function EventHostPanel({
  eventAddress = EVENT_ADDRESS,
  accessManagerAddress = ACCESS_MANAGER_ADDRESS,
}: EventHostPanelProps) {
  const [participantAddress, setParticipantAddress] = useState<Address>(defaultParticipant);
  const [issuerAddress, setIssuerAddress] = useState<Address>(defaultParticipant);
  const [status, setStatus] = useState<EventStatus>(EventStatus.ACTIVE);
  const [eventName, setEventName] = useState("");
  const [eventDescription, setEventDescription] = useState("");
  const [seatsCount, setSeatsCount] = useState<number>(100);
  const [messageHash, setMessageHash] = useState<Hex>("0x");
  const [signature, setSignature] = useState<Hex>("0x");
  const { address } = useAccount();
  const publicClient = usePublicClient({ chainId });
  const { data: walletClient } = useWalletClient({ chainId: sepolia.id });
  const { refreshRoles } = useWalletContext();
  const eventQuery = useEventData(eventAddress);
  const participantMessage = useParticipantMessage(participantAddress, eventAddress);

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

  const handleGrantParticipantRole = async () => {
    if (!participantAddress || participantAddress === defaultParticipant) {
      alert("Enter a participant address");
      return;
    }
    await write({
      address: accessManagerAddress,
      abi: EventAccessManagerABI,
      functionName: "grantParticipantRole",
      args: [participantAddress],
    });
    await refreshRoles();
  };

  const handleRevokeParticipantRole = async () => {
    if (!participantAddress || participantAddress === defaultParticipant) {
      alert("Enter a participant address");
      return;
    }
    await write({
      address: accessManagerAddress,
      abi: EventAccessManagerABI,
      functionName: "revokeParticipantRole",
      args: [participantAddress],
    });
    await refreshRoles();
  };

  const handleGrantIssuerRole = async () => {
    if (!issuerAddress || issuerAddress === defaultParticipant) {
      alert("Enter an issuer address");
      return;
    }
    await write({
      address: accessManagerAddress,
      abi: EventAccessManagerABI,
      functionName: "grantIssuerRole",
      args: [issuerAddress],
    });
  };

  const handleRevokeIssuerRole = async () => {
    if (!issuerAddress || issuerAddress === defaultParticipant) {
      alert("Enter an issuer address");
      return;
    }
    await write({
      address: accessManagerAddress,
      abi: EventAccessManagerABI,
      functionName: "revokeIssuerRole",
      args: [issuerAddress],
    });
  };

  const handleUpdateEvent = async () => {
    await write({
      address: eventAddress,
      abi: EventABI,
      functionName: "updateEvent",
      args: [eventName, eventDescription, BigInt(seatsCount), status],
    });
    await eventQuery.refetch();
  };

  const handleAddParticipant = async () => {
    if (!participantAddress || participantAddress === defaultParticipant) {
      alert("Enter a participant address");
      return;
    }
    await write({
      address: eventAddress,
      abi: EventABI,
      functionName: "addParticipant",
      args: [participantAddress],
    });
  };

  const handleConfirmEvent = async () => {
    await write({
      address: eventAddress,
      abi: EventABI,
      functionName: "confirmEvent",
      args: [],
    });
    await eventQuery.refetch();
  };

  const handleSetSigningMessage = async () => {
    if (!participantAddress || participantAddress === defaultParticipant || !isHex(messageHash)) {
      alert("Provide a participant and valid hash");
      return;
    }
    await write({
      address: eventAddress,
      abi: EventABI,
      functionName: "setSigningMessage",
      args: [participantAddress, messageHash],
    });
    await participantMessage.refetch();
  };

  const handleVerifySignature = async () => {
    if (
      !participantAddress ||
      participantAddress === defaultParticipant ||
      !signature ||
      !isHex(signature)
    ) {
      alert("Provide participant and signature");
      return;
    }
    const recovered = await publicClient?.readContract({
      address: eventAddress,
      abi: EventABI,
      functionName: "verifySignature",
      args: [participantAddress, signature],
    });
    alert(`Recovered address: ${recovered}`);
  };

  const handleMintTicket = async () => {
    if (!EVENT_TICKET_ADDRESS) {
      alert("Ticket contract not configured");
      return;
    }
    await write({
      address: EVENT_TICKET_ADDRESS,
      abi: EventTicketABI,
      functionName: "mintNft",
      args: [participantAddress, "userId", "ticketId", "issuerId", "encryptedData", "backendData"],
    });
  };

  const handleMintCertificate = async () => {
    if (!EVENT_CERTIFICATE_ADDRESS) {
      alert("Certificate contract not configured");
      return;
    }
    await write({
      address: EVENT_CERTIFICATE_ADDRESS,
      abi: EventCertificateABI,
      functionName: "mintNft",
      args: [
        participantAddress,
        "userId",
        "certificateId",
        "issuerId",
        "encryptedData",
        "backendData",
      ],
    });
  };

  return (
    <div className="overflow-hidden rounded-3xl border border-slate-800/60 bg-slate-950/70 shadow-2xl shadow-slate-950/40">
      <div className="relative">
        <div className="absolute inset-0 bg-gradient-to-br from-slate-950 via-emerald-950/25 to-slate-950" />
        <div className="relative grid gap-10 px-8 py-10">
          <div className="flex flex-col gap-3">
            <span className="inline-flex w-fit items-center gap-2 rounded-full border border-emerald-500/30 bg-emerald-500/15 px-3 py-1 text-xs text-emerald-200">
              <span className="inline-flex h-1.5 w-1.5 rounded-full bg-emerald-300" /> Host Panel
            </span>
            <h2 className="text-2xl font-semibold tracking-tight text-slate-100">Host Controls</h2>
            <p className="text-sm text-slate-400">
              Manage roles, confirm attendance, and mint verifiable tickets or certificates for
              approved participants.
            </p>
            {eventQuery.data && (
              <div className="mt-4 grid gap-3 rounded-2xl border border-slate-800/60 bg-slate-950/80 p-5 text-sm text-slate-300">
                <p>
                  <span className="font-semibold text-slate-100">Current Event:</span>{" "}
                  {eventQuery.data.name}
                </p>
                <p className="leading-relaxed text-slate-400">{eventQuery.data.description}</p>
                <div className="flex flex-wrap items-center gap-3 text-xs text-slate-400">
                  <span>
                    Seats {eventQuery.data.currentSeatsCount.toString()} /{" "}
                    {eventQuery.data.seatsCount.toString()}
                  </span>
                  <span className="inline-flex items-center gap-2 rounded-full border border-slate-800/60 bg-slate-900/80 px-3 py-1 text-xs text-slate-300">
                    Status {EventStatus[eventQuery.data.status]}
                  </span>
                </div>
              </div>
            )}
          </div>

          <section className="grid gap-4 lg:grid-cols-2">
            <div className="grid gap-3 rounded-2xl border border-slate-800/60 bg-slate-950/70 p-5">
              <h3 className="text-sm font-medium text-slate-200">Participant Role</h3>
              <label className="text-xs uppercase tracking-widest text-slate-500">
                Wallet Address
              </label>
              <input
                className="rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100 focus:border-emerald-400/50 focus:outline-none focus:ring-2 focus:ring-emerald-400/20"
                value={participantAddress}
                onChange={(event) => setParticipantAddress(event.target.value as Address)}
                placeholder="0xabc..."
              />
              <div className="flex flex-wrap gap-3">
                <button
                  className="rounded-lg bg-emerald-500/15 px-4 py-2 text-sm font-medium text-emerald-200 transition hover:bg-emerald-500/25"
                  onClick={handleGrantParticipantRole}
                >
                  Grant Participant
                </button>
                <button
                  className="rounded-lg bg-rose-500/10 px-4 py-2 text-sm font-medium text-rose-200 transition hover:bg-rose-500/20"
                  onClick={handleRevokeParticipantRole}
                >
                  Revoke Participant
                </button>
              </div>
            </div>

            <div className="grid gap-3 rounded-2xl border border-slate-800/60 bg-slate-950/70 p-5">
              <h3 className="text-sm font-medium text-slate-200">Issuer Role</h3>
              <label className="text-xs uppercase tracking-widest text-slate-500">
                Wallet Address
              </label>
              <input
                className="rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100 focus:border-indigo-400/50 focus:outline-none focus:ring-2 focus:ring-indigo-400/20"
                value={issuerAddress}
                onChange={(event) => setIssuerAddress(event.target.value as Address)}
                placeholder="0xdef..."
              />
              <div className="flex flex-wrap gap-3">
                <button
                  className="rounded-lg bg-indigo-500/15 px-4 py-2 text-sm font-medium text-indigo-200 transition hover:bg-indigo-500/25"
                  onClick={handleGrantIssuerRole}
                >
                  Grant Issuer
                </button>
                <button
                  className="rounded-lg bg-amber-500/10 px-4 py-2 text-sm font-medium text-amber-200 transition hover:bg-amber-500/20"
                  onClick={handleRevokeIssuerRole}
                >
                  Revoke Issuer
                </button>
              </div>
            </div>
          </section>

          <section className="grid gap-5 rounded-2xl border border-slate-800/60 bg-slate-950/70 p-5">
            <h3 className="text-sm font-medium text-slate-200">Event Metadata</h3>
            <div className="grid gap-4 md:grid-cols-2">
              <div className="grid gap-2">
                <label className="text-xs uppercase tracking-widest text-slate-500">
                  Event Name
                </label>
                <input
                  className="rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100 focus:border-sky-400/50 focus:outline-none focus:ring-2 focus:ring-sky-400/20"
                  value={eventName}
                  onChange={(event) => setEventName(event.target.value)}
                  placeholder="Enter new event name"
                />
              </div>
              <div className="grid gap-2">
                <label className="text-xs uppercase tracking-widest text-slate-500">Seats</label>
                <input
                  className="rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100 focus:border-sky-400/50 focus:outline-none focus:ring-2 focus:ring-sky-400/20"
                  type="number"
                  min={0}
                  value={seatsCount}
                  onChange={(event) => setSeatsCount(Number(event.target.value))}
                />
              </div>
            </div>
            <div className="grid gap-2">
              <label className="text-xs uppercase tracking-widest text-slate-500">
                Description
              </label>
              <textarea
                className="rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100 focus:border-sky-400/50 focus:outline-none focus:ring-2 focus:ring-sky-400/20"
                rows={3}
                value={eventDescription}
                onChange={(event) => setEventDescription(event.target.value)}
                placeholder="Describe your event logistics"
              />
            </div>
            <div className="grid gap-3 sm:grid-cols-2">
              <div className="grid gap-2">
                <label className="text-xs uppercase tracking-widest text-slate-500">Status</label>
                <select
                  className="rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100 focus:border-sky-400/50 focus:outline-none focus:ring-2 focus:ring-sky-400/20"
                  value={status}
                  onChange={(event) => setStatus(Number(event.target.value) as EventStatus)}
                >
                  <option value={EventStatus.ACTIVE}>Active</option>
                  <option value={EventStatus.INACTIVE}>Inactive</option>
                  <option value={EventStatus.CLOSED}>Closed</option>
                </select>
              </div>
              <div className="flex items-end justify-end gap-3">
                <button
                  className="inline-flex w-full justify-center rounded-lg bg-sky-500/15 px-4 py-2 text-sm font-medium text-sky-200 transition hover:bg-sky-500/25 sm:w-auto"
                  onClick={handleAddParticipant}
                >
                  Add Participant
                </button>
                <button
                  className="inline-flex w-full justify-center rounded-lg bg-purple-500/15 px-4 py-2 text-sm font-medium text-purple-200 transition hover:bg-purple-500/25 sm:w-auto"
                  onClick={handleConfirmEvent}
                >
                  Confirm Event
                </button>
              </div>
            </div>
            <button
              className="rounded-lg border border-sky-400/40 bg-sky-500/15 px-4 py-2 text-sm font-medium text-sky-200 transition hover:bg-sky-500/25"
              onClick={handleUpdateEvent}
            >
              Save Event Changes
            </button>
          </section>

          <section className="grid gap-4 rounded-2xl border border-slate-800/60 bg-slate-950/70 p-5">
            <h3 className="text-sm font-medium text-slate-200">Signing Workflow</h3>
            <div className="grid gap-3 text-sm">
              <div className="flex flex-col gap-2">
                <span className="text-xs uppercase tracking-widest text-slate-500">
                  Current Message
                </span>
                <code className="rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-xs text-slate-300">
                  {participantMessage.data ?? "0x"}
                </code>
              </div>
              <div className="grid gap-2">
                <label className="text-xs uppercase tracking-widest text-slate-500">
                  Message Hash
                </label>
                <input
                  className="rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100 focus:border-indigo-400/50 focus:outline-none focus:ring-2 focus:ring-indigo-400/20"
                  value={messageHash}
                  onChange={(event) => setMessageHash(event.target.value as Hex)}
                  placeholder="0x message hash"
                />
              </div>
              <div className="grid gap-2">
                <label className="text-xs uppercase tracking-widest text-slate-500">
                  Participant Signature
                </label>
                <input
                  className="rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100 focus:border-amber-400/50 focus:outline-none focus:ring-2 focus:ring-amber-400/20"
                  value={signature}
                  onChange={(event) => setSignature(event.target.value as Hex)}
                  placeholder="0x signature"
                />
              </div>
              <div className="flex flex-wrap gap-3">
                <button
                  className="rounded-lg bg-indigo-500/15 px-4 py-2 text-sm font-medium text-indigo-200 transition hover:bg-indigo-500/25"
                  onClick={handleSetSigningMessage}
                >
                  Push Signing Message
                </button>
                <button
                  className="rounded-lg bg-amber-500/15 px-4 py-2 text-sm font-medium text-amber-200 transition hover:bg-amber-500/25"
                  onClick={handleVerifySignature}
                >
                  Verify Signature
                </button>
              </div>
            </div>
          </section>

          <section className="grid gap-3 rounded-2xl border border-slate-800/60 bg-slate-950/70 p-5">
            <h3 className="text-sm font-medium text-slate-200">Credential Issuance</h3>
            <div className="flex flex-wrap gap-3">
              <button
                className="rounded-lg bg-emerald-500/15 px-4 py-2 text-sm font-medium text-emerald-200 transition hover:bg-emerald-500/25"
                onClick={handleMintTicket}
              >
                Mint Ticket NFT
              </button>
              <button
                className="rounded-lg bg-violet-500/15 px-4 py-2 text-sm font-medium text-violet-200 transition hover:bg-violet-500/25"
                onClick={handleMintCertificate}
              >
                Mint Certificate NFT
              </button>
            </div>
            <p className="text-xs text-slate-400">
              Ticket contract:{" "}
              <code className="text-slate-300">{EVENT_TICKET_ADDRESS || "not configured"}</code>
              <br />
              Certificate contract:{" "}
              <code className="text-slate-300">
                {EVENT_CERTIFICATE_ADDRESS || "not configured"}
              </code>
            </p>
          </section>
        </div>
      </div>
    </div>
  );
}
