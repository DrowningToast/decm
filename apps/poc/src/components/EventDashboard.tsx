import { useAccount } from "wagmi";
import { EventHostPanel } from "./EventHostPanel";
import { EventParticipantPanel } from "./EventParticipantPanel";
import { useWalletContext } from "../context/WalletContext";
import { ACCESS_MANAGER_ADDRESS, EVENT_ADDRESS } from "../lib/addresses";
import { useEventData } from "../hooks/useEventData";
import { EventStatus } from "../types/Event";
import { cn } from "../lib/utils";

function renderStatusLabel(status: EventStatus) {
  switch (status) {
    case EventStatus.ACTIVE:
      return {
        label: "Active",
        className: "bg-emerald-500/15 text-emerald-300",
      };
    case EventStatus.INACTIVE:
      return { label: "Inactive", className: "bg-amber-500/15 text-amber-300" };
    case EventStatus.CLOSED:
      return { label: "Closed", className: "bg-rose-500/15 text-rose-300" };
    default:
      return {
        label: `Unknown (${status})`,
        className: "bg-slate-500/15 text-slate-300",
      };
  }
}

export function EventDashboard() {
  const { address } = useAccount();
  const { isHost, isParticipant, isIssuer } = useWalletContext();
  const eventQuery = useEventData(EVENT_ADDRESS);

  if (!address) {
    return (
      <div className="rounded-2xl border border-slate-800/80 bg-slate-900/80 px-8 py-16 text-center shadow-lg shadow-slate-900/30">
        <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full border border-slate-800 bg-slate-900 text-slate-400">
          <svg viewBox="0 0 24 24" className="h-8 w-8" aria-hidden>
            <path
              fill="currentColor"
              d="M12 2a7 7 0 0 0-7 7v3.278a2 2 0 0 1-.553 1.381l-.724.761A2 2 0 0 0 5.724 18h12.552a2 2 0 0 0 1.001-3.58l-.724-.761A2 2 0 0 1 18 12.278V9a7 7 0 0 0-7-7Zm0 20a3 3 0 0 0 3-3H9a3 3 0 0 0 3 3Z"
            />
          </svg>
        </div>
        <h2 className="mt-6 text-2xl font-semibold tracking-tight text-slate-100">
          Wallet Connection Required
        </h2>
        <p className="mt-2 text-sm text-slate-400">
          Connect MetaMask to view event details, manage roles, and fetch on-chain data from
          Sepolia.
        </p>
      </div>
    );
  }

  return (
    <div className="grid gap-10">
      {eventQuery.data && (
        <section className="overflow-hidden rounded-3xl border border-slate-800/60 bg-slate-950/60 shadow-2xl shadow-slate-950/40">
          <div className="relative">
            <div className="absolute inset-0 bg-gradient-to-r from-emerald-500/5 via-slate-950 to-blue-500/5" />
            <div className="relative grid gap-6 px-8 py-8 sm:grid-cols-[1.2fr,0.8fr]">
              <div>
                <p className="text-xs uppercase tracking-[0.3em] text-slate-400">Current Event</p>
                <h2 className="mt-3 text-3xl font-semibold text-slate-50">
                  {eventQuery.data.name}
                </h2>
                <p className="mt-3 text-sm leading-relaxed text-slate-300">
                  {eventQuery.data.description}
                </p>
              </div>
              <div className="flex flex-col justify-between">
                <div className="flex flex-wrap items-center gap-4">
                  <span className="inline-flex items-center gap-2 rounded-full border border-slate-800/50 bg-slate-900/80 px-4 py-2 text-sm text-slate-300">
                    <svg viewBox="0 0 24 24" className="h-4 w-4 text-emerald-300" aria-hidden>
                      <path
                        fill="currentColor"
                        d="M19.5 5h-15A2.5 2.5 0 0 0 2 7.5v9A2.5 2.5 0 0 0 4.5 19h15a2.5 2.5 0 0 0 2.5-2.5v-9A2.5 2.5 0 0 0 19.5 5Zm-15-1.5h15A4 4 0 0 1 23 7.5v9a4 4 0 0 1-4 4h-15a4 4 0 0 1-4-4v-9a4 4 0 0 1 4-4Z"
                      />
                    </svg>
                    Seats {eventQuery.data.currentSeatsCount.toString()} /{" "}
                    {eventQuery.data.seatsCount.toString()}
                  </span>
                  {(() => {
                    const status = renderStatusLabel(eventQuery.data.status);
                    return (
                      <span
                        className={cn(
                          "inline-flex items-center gap-2 rounded-full px-4 py-2 text-sm",
                          status.className
                        )}
                      >
                        <span className="inline-flex h-2 w-2 rounded-full bg-current" />
                        {status.label}
                      </span>
                    );
                  })()}
                </div>
                <div className="mt-6 flex flex-col gap-2 text-xs text-slate-500">
                  <span>
                    Event Contract: <code className="text-slate-300">{EVENT_ADDRESS}</code>
                  </span>
                  <span>
                    Access Manager: <code className="text-slate-300">{ACCESS_MANAGER_ADDRESS}</code>
                  </span>
                </div>
              </div>
            </div>
          </div>
        </section>
      )}

      <div className="grid gap-6">
        {(isHost || isIssuer) && (
          <EventHostPanel
            eventAddress={EVENT_ADDRESS}
            accessManagerAddress={ACCESS_MANAGER_ADDRESS}
          />
        )}

        {!(isHost || isIssuer) && <EventParticipantPanel eventAddress={EVENT_ADDRESS} />}

        {!isHost && !isIssuer && !isParticipant && (
          <div className="rounded-2xl border border-slate-800/80 bg-slate-950/80 px-6 py-10 text-center text-slate-400">
            <h3 className="text-lg font-medium text-slate-100">No Event Roles Assigned</h3>
            <p className="mt-2 text-sm">
              Ask the event host to grant you participant or issuer access via the control panel
              above.
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
