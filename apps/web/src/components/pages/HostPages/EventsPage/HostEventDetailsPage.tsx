import PageContainer from "@/components/container/PageContainer";
import SectionContainer from "@/components/container/SectionContainer";
import { Typography } from "@/components/typography/typography";
import { GoogleMapsEmbed } from "@/components/ui/google-maps-embed";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
} from "@/components/ui/accordion";
import { RequirementItem, type RequirementStatus } from "@/components/ui/requirement-item";
import { TextLabelValue } from "@/components/ui/text-label-value";
import WrappedButton from "@/components/wrapper/WrappedButton";
import { CheckCircle2Icon, CogIcon, ExternalLinkIcon } from "lucide-react";
import { useTranslation } from "react-i18next";

interface HostEventDetailsPageProps {
    eventId: string;
}

export default function HostEventDetailsPage({ eventId }: HostEventDetailsPageProps) {
    const { t } = useTranslation();

    // Mock participant requirements data - replace with actual data
    const participantRequirements: Record<string, RequirementStatus> = {
        firstName: "required",
        lastName: "required",
        email: "required",
        bio: "optional",
        phoneNumber: "optional",
        address: "not_required",
        academicInstitution: "required",
        academicEmail: "required",
    };

    const eventSettings = {
        eventType: "Public",
        bookingRequired: true,
        tokenTransferable: false,
    };

    const MockTable = (
        <div className="overflow-x-auto">
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
                            <button className="text-primary underline hover:opacity-80 transition">
                                View
                            </button>
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
                            <button className="text-primary underline hover:opacity-80 transition">
                                View
                            </button>
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
                    <div className="flex items-center gap-4">
                        <div className="w-8 h-8 bg-white rounded-full"></div>
                        <div className="flex items-center gap-2">
                            <Typography tag="p" size={"subheader"}>
                                ToBeIT 67
                            </Typography>
                            <CheckCircle2Icon color="#eb5331" />
                        </div>
                    </div>

                    <WrappedButton href={`/host/events/${eventId}/edit`}>Edit Event</WrappedButton>
                </div>
            </SectionContainer>

            <SectionContainer className="lg:grid lg:grid-cols-4 gap-8">
                <div className="flex flex-col gap-4 lg:col-span-3">
                    <Typography tag="p" size={"base"} color="muted">
                        Lorem ipsum dolor sit amet consectetur adipisicing elit. Fuga magni
                        exercitationem debitis magnam accusamus ipsum! Velit, provident repellendus
                        quas tempora, odio, quae iure repudiandae est rerum ea quibusdam maiores!
                        In?
                    </Typography>

                    <Typography tag="p" size={"base"} color="muted">
                        Lorem ipsum dolor sit amet consectetur adipisicing elit. Fuga magni
                        exercitationem debitis magnam accusamus ipsum! Velit, provident repellendus
                        quas tempora, odio, quae iure repudiandae est rerum ea quibusdam maiores!
                        In?
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
                <Tabs defaultValue="event-info">
                    <TabsList className="w-full h-10 bg-[#E9DEDE]">
                        <TabsTrigger
                            value="event-info"
                            className="data-[state=active]:bg-primary text-gray-900 data-[state=active]:text-white"
                            style={{
                                fontFamily: "Cormorant Garamond",
                                fontSize: "16px",
                            }}
                        >
                            Event Info
                        </TabsTrigger>
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
                    <TabsContent value="event-info">
                        <div className="mt-6 flex flex-col gap-4 lg:flex-row">
                            <div className="flex flex-col gap-4 flex-1">
                                <TextLabelValue
                                    label="Venue Location"
                                    value="School of Information Technology, KMITL"
                                />
                                <TextLabelValue
                                    label="Google Map Search"
                                    value="School of Information Technology, KMITL"
                                />

                                <TextLabelValue
                                    label="Contract Address"
                                    value="Bangkok, Thailand"
                                />
                                <TextLabelValue label="Contact" value="0656526769" />
                            </div>
                            <div className="flex-1">
                                <GoogleMapsEmbed query="School of Information Technology, KMITL" />
                            </div>
                        </div>
                    </TabsContent>
                    <TabsContent value="participants">
                        <div className="mt-6 space-y-4">
                            {/* Event Settings Summary */}
                            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 my-6">
                                <TextLabelValue
                                    label={t("events.settings.eventType")}
                                    value={eventSettings.eventType}
                                />
                                <TextLabelValue
                                    label={t("events.settings.bookingRequired")}
                                    value={
                                        eventSettings.bookingRequired
                                            ? t("common.yes")
                                            : t("common.no")
                                    }
                                />
                                <TextLabelValue
                                    label={t("events.settings.tokenTransferable")}
                                    value={
                                        eventSettings.tokenTransferable
                                            ? t("common.yes")
                                            : t("common.no")
                                    }
                                />

                                <WrappedButton>
                                    <CogIcon />
                                    {t("events.settings.participantSettings")}
                                </WrappedButton>
                            </div>

                            {/* Participant Requirements Accordion */}
                            <Accordion
                                type="single"
                                collapsible
                                className="w-full bg-white rounded-lg "
                            >
                                <AccordionItem value="requirements">
                                    <AccordionTrigger className="hover:no-underline px-4">
                                        <div className="flex items-center justify-between w-full">
                                            <Typography
                                                variant="text"
                                                tag="span"
                                                className="font-medium text-black"
                                            >
                                                {t("events.participants.requirementsTitle")}
                                            </Typography>
                                        </div>
                                    </AccordionTrigger>
                                    <AccordionContent>
                                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 p-4">
                                            <RequirementItem
                                                label={t("events.participants.fields.firstName")}
                                                status={participantRequirements.firstName}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.lastName")}
                                                status={participantRequirements.lastName}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.email")}
                                                status={participantRequirements.email}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.bio")}
                                                status={participantRequirements.bio}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.phoneNumber")}
                                                status={participantRequirements.phoneNumber}
                                            />
                                            <RequirementItem
                                                label={t("events.participants.fields.address")}
                                                status={participantRequirements.address}
                                            />
                                            <RequirementItem
                                                label={t(
                                                    "events.participants.fields.academicInstitution",
                                                )}
                                                status={participantRequirements.academicInstitution}
                                            />
                                            <RequirementItem
                                                label={t(
                                                    "events.participants.fields.academicEmail",
                                                )}
                                                status={participantRequirements.academicEmail}
                                            />
                                        </div>
                                    </AccordionContent>
                                </AccordionItem>
                            </Accordion>

                            {MockTable}
                        </div>
                    </TabsContent>
                    <TabsContent value="certificates">
                        <div className="mt-6 w-full bg-primary/10 border border-primary/20 rounded-lg p-6 flex flex-col items-center justify-center mb-6">
                            <p className="font-semibold text-lg text-primary">
                                Add event&apos;s certificate configuration
                            </p>
                            <p className="text-muted-foreground text-base mt-1 text-center max-w-xl">
                                Set up the certificate template and rules for this event.
                                Participants will receive certificates based on your configuration.
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
