import PageContainer from "@/components/container/PageContainer";
import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";
import { Typography } from "@/components/typography/typography";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import WrappedButton from "@/components/wrapper/WrappedButton";
import { cn } from "@/lib/utils";
import { ExternalLinkIcon } from "lucide-react";

interface HostEventDetailsPageProps {
  eventId: string;
}

export default function HostEventDetailsPage({ eventId }: HostEventDetailsPageProps) {
  const MockTable = (
    <div className="overflow-x-auto mt-6">
      <table className="min-w-full rounded bg-white text-black">
        <thead className=" h-12 text-black">
          <tr>
            <th className="px-4 py-2 text-left font-medium">Name</th>
            <th className="px-4 py-2 text-left font-medium">Email</th>
            <th className="px-4 py-2 text-left font-medium">Phone Number</th>
            <th className="px-4 py-2 text-left font-medium">Wallet Address</th>
            <th className="px-4 py-2 text-left font-medium">Status</th>
            <th className="px-4 py-2 text-left font-medium">Action</th>
          </tr>
        </thead>
        <tbody>
          {/* Example row, replace with dynamic data as needed */}
          <tr className="border-b">
            <td className="px-4 py-2">John Doe</td>
            <td className="px-4 py-2">john.doe@email.com</td>
            <td className="px-4 py-2">+1 234 567 8901</td>
            <td className="px-4 py-2">0x0000...0000</td>
            <td className="px-4 py-2">
              <span className="inline-block px-2 py-1 rounded bg-green-100 text-green-800 text-xs">
                Confirmed
              </span>
            </td>
            <td className="px-4 py-2">
              <button className="text-primary underline hover:opacity-80 transition">View</button>
            </td>
          </tr>
          <tr>
            <td className="px-4 py-2">Jane Smith</td>
            <td className="px-4 py-2">jane.smith@email.com</td>
            <td className="px-4 py-2">+1 987 654 3210</td>
            <td className="px-4 py-2">0x0000...0000</td>
            <td className="px-4 py-2">
              <span className="inline-block px-2 py-1 rounded bg-yellow-100 text-yellow-800 text-xs">
                Pending
              </span>
            </td>
            <td className="px-4 py-2">
              <button className="text-primary underline hover:opacity-80 transition">View</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );

  return (
    <PageContainer title="Events Details">
      <SectionContainer>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-white rounded-full"></div>
            <Typography tag="p" size={"subheader"}>
              ToBeIT 67
            </Typography>
          </div>

          <WrappedButton href={`/host/events/${eventId}/edit`}>Edit Event</WrappedButton>
        </div>
      </SectionContainer>

      <SectionContainer className="lg:grid lg:grid-cols-4 gap-8">
        <div className="flex flex-col gap-4 lg:col-span-3">
          <Typography tag="p" size={"base"} color="muted">
            Lorem ipsum dolor sit amet consectetur adipisicing elit. Fuga magni exercitationem
            debitis magnam accusamus ipsum! Velit, provident repellendus quas tempora, odio, quae
            iure repudiandae est rerum ea quibusdam maiores! In?
          </Typography>

          <div className="w-full h-[175px] bg-white rounded mt-1"></div>
        </div>

        <div className="flex flex-col gap-4 mt-6 lg:mt-0">
          <TextLabelValue label="Status" value="ToBeIT 67" />
          <TextLabelValue label="Final call for request" value="2025-01-01" />
          <TextLabelValue label="Participation request" value="Invited Only" />
          <TextLabelValue label="Seats count" value="20/40" />
          <TextLabelValue
            label="Event Contract Address"
            value="0x0000...0000"
            endIcon={<ExternalLinkIcon size={16} />}
            valueClassName="cursor-pointer underline"
            href="https://www.google.com"
          />
        </div>
      </SectionContainer>

      <SectionContainer>
        <Tabs defaultValue="participants">
          <TabsList className="w-full h-10 bg-[#E9DEDE]">
            <TabsTrigger
              value="participants"
              className="data-[state=active]:bg-primary text-gray-900 data-[state=active]:text-white"
              style={{
                fontFamily: "Cormorant Garamond",
                fontSize: "16px",
              }}
            >
              Participants
            </TabsTrigger>
            <TabsTrigger
              value="certificates"
              className="data-[state=active]:bg-primary text-gray-900 data-[state=active]:text-white"
              style={{
                fontFamily: "Cormorant Garamond",
                fontSize: "16px",
              }}
            >
              Certificates
            </TabsTrigger>
          </TabsList>
          <TabsContent value="participants">{MockTable}</TabsContent>
          <TabsContent value="certificates">
            <div className="mt-6 w-full bg-primary/10 border border-primary/20 rounded-lg p-6 flex flex-col items-center justify-center mb-6">
              <p className="font-semibold text-lg text-primary">
                Add event&apos;s certificate configuration
              </p>
              <p className="text-muted-foreground text-base mt-1 text-center max-w-xl">
                Set up the certificate template and rules for this event. Participants will receive
                certificates based on your configuration.
              </p>
              <button
                className="mt-4 px-5 py-2 rounded-md bg-primary text-white font-medium hover:bg-primary/90 transition"
                type="button"
              >
                Configure Certificates
              </button>
            </div>

            {MockTable}
          </TabsContent>
        </Tabs>
      </SectionContainer>
    </PageContainer>
  );
}

interface TextLabelValueProps {
  label: string;
  value: string;
  endIcon?: React.ReactNode;
  valueClassName?: string;
  href?: string;
}
function TextLabelValue({ label, value, endIcon, valueClassName, href }: TextLabelValueProps) {
  return (
    <div className="flex flex-col gap-1">
      <Typography tag="p" size={"base"} color="muted" className="text-sm">
        {label}
      </Typography>
      <Typography
        tag="p"
        size={"base"}
        className={cn(valueClassName, {
          "flex items-center gap-2": endIcon,
        })}
      >
        {href ? (
          <a href={href} target="_blank">
            {value}
          </a>
        ) : (
          value
        )}
        {endIcon}
      </Typography>
    </div>
  );
}
