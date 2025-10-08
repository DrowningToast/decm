import PageContainer from "@/components/container/PageContainer";
import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";
import WrappedButton from "@/components/wrapper/WrappedButton";
import { Link } from "@/router";

export default function HostEventPage() {
  return (
    <PageContainer title="Events">
      <SectionContainer className="flex items-center justify-between">
        <TitleSubtitle title="Events" subtitle="Create or manage your events" />
        <div className="flex justify-end">
          <WrappedButton href="/host/events/create">Create Event</WrappedButton>
        </div>
      </SectionContainer>

      <SectionContainer>
        <table className="w-full">
          <thead className="border-b  h-10">
            <tr>
              <th className="text-start text-muted">Name</th>
              <th className="text-end text-muted">Status</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>
                <Link
                  to={"/host/events/:eventId"}
                  params={{
                    eventId: "1",
                  }}
                >
                  <p className="underline mt-4 cursor-pointer">ToBeIT 67</p>
                </Link>
                <p className="text-muted text-sm mt-0.5">24 Sep 2025 - 25 Sep 2025</p>
                <p className="text-muted text-sm">Invite Only</p>
              </td>
              <td className="text-end">Active</td>
            </tr>

            <tr>
              <td>
                <p className="underline mt-4">ToBeIT 67</p>
                <p className="text-muted text-sm mt-0.5">24 Sep 2025 - 25 Sep 2025</p>
                <p className="text-muted text-sm">Invite Only</p>
              </td>
              <td className="text-end">Active</td>
            </tr>
          </tbody>
        </table>
      </SectionContainer>
    </PageContainer>
  );
}
