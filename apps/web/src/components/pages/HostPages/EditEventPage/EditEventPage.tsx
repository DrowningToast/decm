import { useTranslation } from "react-i18next";
import { EventForm } from "@/components/forms/EventForm";
import type { EventFormData } from "@/lib/schemas/eventFormSchema";
import { toast } from "sonner";
import PageContainer from "@/components/container/PageContainer";
import SectionContainer from "@/components/container/SectionContainer";
import TitleSubtitle from "@/components/TitleSubtitle";

export const EditEventPage = () => {
  const { t } = useTranslation();

  const handleEditEvent = async (data: EventFormData) => {
    try {
      // TODO: Implement API call to edit event
      console.log("Editing event:", data);

      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 1000));

      // Show success toast
      toast.success(t("common.success"), {
        description: `Event "${data.name}" edited successfully`,
      });
    } catch (error) {
      // Show error toast
      toast.error(t("common.error"), {
        description: t("errors.generic"),
      });
      console.error("Error editing event:", error);
    }
  };

  return (
    <PageContainer title="Edit Event" className="space-y-6">
      {/* Page Header */}
      <SectionContainer>
        <TitleSubtitle title="Edit Event" subtitle="Fill in the details below to edit the event" />
      </SectionContainer>

      <SectionContainer>
        <EventForm onSubmit={handleEditEvent} mode="edit" />
      </SectionContainer>
    </PageContainer>
  );
};
