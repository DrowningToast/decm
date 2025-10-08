import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslation } from "react-i18next";
import { useState } from "react";
import { format } from "date-fns";
import { Calendar as CalendarIcon, Upload, X } from "lucide-react";
import { eventFormSchema, type EventFormData } from "@/lib/schemas/eventFormSchema";
import { Typography } from "@/components/typography/typography";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { cn } from "@/lib/utils";

interface EventFormProps {
  /**
   * Default values for the form (used for editing existing events)
   */
  defaultValues?: Partial<EventFormData>;
  /**
   * Callback function called when form is successfully submitted
   */
  onSubmit: (data: EventFormData) => void | Promise<void>;
  /**
   * Whether the form is in loading/submitting state
   */
  isLoading?: boolean;
  /**
   * Form mode: 'create' or 'edit'
   */
  mode?: "create" | "edit";
}

export const EventForm = ({
  defaultValues,
  onSubmit,
  isLoading = false,
  mode = "create",
}: EventFormProps) => {
  const { t } = useTranslation();

  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [iconPreview, setIconPreview] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    setValue,
    watch,
  } = useForm<EventFormData>({
    resolver: zodResolver(eventFormSchema),
    defaultValues: {
      name: defaultValues?.name || "",
      description: defaultValues?.description || "",
      eventBanner: defaultValues?.eventBanner || undefined,
      eventIcon: defaultValues?.eventIcon || undefined,
      startDate: defaultValues?.startDate || undefined,
      endDate: defaultValues?.endDate || undefined,
      seatsCount: defaultValues?.seatsCount || undefined,
    },
  });

  const eventBanner = watch("eventBanner");
  const eventIcon = watch("eventIcon");

  const handleFormSubmit = async (data: EventFormData) => {
    await onSubmit(data);
  };

  const handleImageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setValue("eventBanner", file);
      const reader = new FileReader();
      reader.onloadend = () => {
        setImagePreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleRemoveImage = () => {
    setValue("eventBanner", undefined);
    setImagePreview(null);
  };

  const handleIconChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setValue("eventIcon", file);
      const reader = new FileReader();
      reader.onloadend = () => {
        setIconPreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleRemoveIcon = () => {
    setValue("eventIcon", undefined);
    setIconPreview(null);
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {/* Event Banner */}
        <div className="space-y-2 lg:col-span-3">
          <Label htmlFor="eventBanner">
            <Typography variant="text" tag="span" className="text-sm font-medium">
              {t("events.form.eventBanner")}
            </Typography>
          </Label>

          {!imagePreview && !eventBanner ? (
            <div className="relative">
              <input
                type="file"
                id="eventBanner"
                accept="image/jpeg,image/jpg,image/png,image/webp"
                className="hidden"
                onChange={handleImageChange}
                disabled={isLoading}
              />
              <label
                htmlFor="eventBanner"
                className={cn(
                  "flex flex-col items-center justify-center w-full h-48 border-2 border-dashed rounded-lg cursor-pointer transition-colors",
                  "border-gray-300 bg-background hover:bg-accent/50",
                  isLoading && "opacity-50 cursor-not-allowed"
                )}
              >
                <div className="flex flex-col items-center justify-center pt-5 pb-6">
                  <Upload className="w-10 h-10 mb-3 text-muted-foreground" />
                  <Typography variant="text" tag="p" className="mb-2 text-sm text-muted-foreground">
                    <span className="font-semibold">{t("events.form.eventBannerPlaceholder")}</span>
                  </Typography>
                  <Typography variant="text" tag="p" className="text-xs text-muted-foreground">
                    PNG, JPG, or WebP (MAX. 5MB)
                  </Typography>
                </div>
              </label>
            </div>
          ) : (
            <div className="space-y-3">
              <div className="relative w-full h-72 rounded-lg overflow-hidden border border-[#D9D9D91A]">
                <img
                  src={imagePreview || ""}
                  alt="Event banner preview"
                  className="w-full h-full object-cover"
                />
              </div>
              <div className="flex gap-2">
                <label htmlFor="eventBanner-change" className="flex-1">
                  <input
                    type="file"
                    id="eventBanner-change"
                    accept="image/jpeg,image/jpg,image/png,image/webp"
                    className="hidden"
                    onChange={handleImageChange}
                    disabled={isLoading}
                  />
                  <Button
                    type="button"
                    variant="secondary-dark"
                    size="default"
                    disabled={isLoading}
                    className="w-full cursor-pointer"
                    asChild
                  >
                    <span>
                      <Upload className="h-4 w-4 mr-2" />
                      <Typography variant="text" tag="span" className="font-medium">
                        {t("events.form.eventBannerChange")}
                      </Typography>
                    </span>
                  </Button>
                </label>
                <Button
                  type="button"
                  variant="secondary-dark"
                  size="default"
                  onClick={handleRemoveImage}
                  disabled={isLoading}
                  className="flex-1"
                >
                  <X className="h-4 w-4 mr-2" />
                  <Typography variant="text" tag="span" className="font-medium">
                    {t("events.form.eventBannerRemove")}
                  </Typography>
                </Button>
              </div>
            </div>
          )}

          {errors.eventBanner && (
            <Typography variant="text" tag="p" className="text-sm text-destructive" role="alert">
              {t(errors.eventBanner.message as string)}
            </Typography>
          )}
        </div>

        {/* Event Icon */}
        <div className="space-y-2">
          <Label htmlFor="eventIcon">
            <Typography variant="text" tag="span" className="text-sm font-medium">
              {t("events.form.eventIcon")}
            </Typography>
          </Label>

          {!iconPreview && !eventIcon ? (
            <div className="relative">
              <input
                type="file"
                id="eventIcon"
                accept="image/jpeg,image/jpg,image/png,image/webp"
                className="hidden"
                onChange={handleIconChange}
                disabled={isLoading}
              />
              <label
                htmlFor="eventIcon"
                className={cn(
                  "flex flex-col items-center justify-center w-full h-48 border-2 border-dashed rounded-lg cursor-pointer transition-colors",
                  "border-input bg-background hover:bg-accent/50",
                  isLoading && "opacity-50 cursor-not-allowed"
                )}
              >
                <div className="flex flex-col items-center justify-center pt-5 pb-6">
                  <Upload className="w-10 h-10 mb-3 text-muted-foreground" />
                  <Typography variant="text" tag="p" className="mb-2 text-sm text-muted-foreground">
                    <span className="font-semibold">{t("events.form.eventIconPlaceholder")}</span>
                  </Typography>
                  <Typography variant="text" tag="p" className="text-xs text-muted-foreground">
                    PNG, JPG, or WebP (MAX. 5MB)
                  </Typography>
                </div>
              </label>
            </div>
          ) : (
            <div className="space-y-3">
              <div className="relative w-full h-72 rounded-lg overflow-hidden border border-[#D9D9D91A]">
                <img src={iconPreview || ""} alt="Event icon preview" className="w-full h-full " />
              </div>
              <div className="flex gap-2">
                <label htmlFor="eventIcon-change" className="flex-1">
                  <input
                    type="file"
                    id="eventIcon-change"
                    accept="image/jpeg,image/jpg,image/png,image/webp"
                    className="hidden"
                    onChange={handleIconChange}
                    disabled={isLoading}
                  />
                  <Button
                    type="button"
                    variant="secondary-dark"
                    size="default"
                    disabled={isLoading}
                    className="w-full cursor-pointer"
                    asChild
                  >
                    <span>
                      <Upload className="h-4 w-4 mr-2" />
                      <Typography variant="text" tag="span" className="font-medium">
                        {t("events.form.eventIconChange")}
                      </Typography>
                    </span>
                  </Button>
                </label>
                <Button
                  type="button"
                  variant="secondary-dark"
                  size="default"
                  onClick={handleRemoveIcon}
                  disabled={isLoading}
                  className="flex-1"
                >
                  <X className="h-4 w-4 mr-2" />
                  <Typography variant="text" tag="span" className="font-medium">
                    {t("events.form.eventIconRemove")}
                  </Typography>
                </Button>
              </div>
            </div>
          )}

          {errors.eventIcon && (
            <Typography variant="text" tag="p" className="text-sm text-destructive" role="alert">
              {t(errors.eventIcon.message as string)}
            </Typography>
          )}
        </div>
      </div>

      {/* Event Name */}
      <div className="space-y-2">
        <Label htmlFor="name">
          <Typography variant="text" tag="span" className="text-sm font-medium">
            {t("events.form.name")}
          </Typography>
          <span className="text-destructive ml-1">*</span>
        </Label>
        <Input
          id="name"
          type="text"
          placeholder={t("events.form.namePlaceholder")}
          aria-invalid={!!errors.name}
          disabled={isLoading}
          className="!border !border-[#D9D9D91A]"
          {...register("name")}
        />
        {errors.name && (
          <Typography variant="text" tag="p" className="text-sm text-destructive" role="alert">
            {t(errors.name.message as string)}
          </Typography>
        )}
      </div>

      {/* Description */}
      <div className="space-y-2">
        <Label htmlFor="description">
          <Typography variant="text" tag="span" className="text-sm font-medium">
            {t("events.form.description")}
          </Typography>
        </Label>
        <Textarea
          id="description"
          placeholder={t("events.form.descriptionPlaceholder")}
          disabled={isLoading}
          className="!border !border-[#D9D9D91A]"
          {...register("description")}
        />
        {errors.description && (
          <Typography variant="text" tag="p" className="text-sm text-destructive" role="alert">
            {t(errors.description.message as string)}
          </Typography>
        )}
      </div>

      {/* Date Fields */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Start Date */}
        <div className="space-y-2">
          <Label htmlFor="startDate">
            <Typography variant="text" tag="span" className="text-sm font-medium">
              {t("events.form.startDate")}
            </Typography>
            <span className="text-destructive ml-1">*</span>
          </Label>
          <Controller
            name="startDate"
            control={control}
            render={({ field }) => (
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                    variant="secondary-dark"
                    id="startDate"
                    className={cn(
                      "w-full justify-start text-left font-normal border border-[#D9D9D91A] bg-transparent",
                      !field.value && "text-muted-foreground"
                    )}
                    disabled={isLoading}
                  >
                    <CalendarIcon className="mr-2 h-4 w-4" />
                    {field.value
                      ? format(field.value, "PPP")
                      : t("events.form.startDatePlaceholder")}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={field.value}
                    onSelect={field.onChange}
                    disabled={(date) => date < new Date(new Date().setHours(0, 0, 0, 0))}
                    initialFocus
                  />
                </PopoverContent>
              </Popover>
            )}
          />
          {errors.startDate && (
            <Typography variant="text" tag="p" className="text-sm text-destructive" role="alert">
              {t(errors.startDate.message as string)}
            </Typography>
          )}
        </div>

        {/* End Date */}
        <div className="space-y-2">
          <Label htmlFor="endDate">
            <Typography variant="text" tag="span" className="text-sm font-medium">
              {t("events.form.endDate")}
            </Typography>
            <span className="text-destructive ml-1">*</span>
          </Label>
          <Controller
            name="endDate"
            control={control}
            render={({ field }) => (
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                    variant="secondary-dark"
                    id="endDate"
                    className={cn(
                      "w-full justify-start text-left font-normal border border-[#D9D9D91A] bg-transparent",
                      !field.value && "text-muted-foreground"
                    )}
                    disabled={isLoading}
                  >
                    <CalendarIcon className="mr-2 h-4 w-4" />
                    {field.value ? format(field.value, "PPP") : t("events.form.endDatePlaceholder")}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={field.value}
                    onSelect={field.onChange}
                    disabled={(date) => date < new Date(new Date().setHours(0, 0, 0, 0))}
                    initialFocus
                  />
                </PopoverContent>
              </Popover>
            )}
          />
          {errors.endDate && (
            <Typography variant="text" tag="p" className="text-sm text-destructive" role="alert">
              {t(errors.endDate.message as string)}
            </Typography>
          )}
        </div>
      </div>

      {/* Seats Count */}
      <div className="space-y-2">
        <Label htmlFor="seatsCount">
          <Typography variant="text" tag="span" className="text-sm font-medium">
            {t("events.form.seatsCount")}
          </Typography>
          <span className="text-destructive ml-1">*</span>
        </Label>
        <Input
          id="seatsCount"
          type="number"
          min="1"
          step="1"
          placeholder={t("events.form.seatsCountPlaceholder")}
          aria-invalid={!!errors.seatsCount}
          disabled={isLoading}
          className="!border !border-[#D9D9D91A]"
          {...register("seatsCount", { valueAsNumber: true })}
        />
        {errors.seatsCount && (
          <Typography variant="text" tag="p" className="text-sm text-destructive" role="alert">
            {t(errors.seatsCount.message as string)}
          </Typography>
        )}
      </div>

      {/* Submit Button */}
      <div className="flex justify-end gap-4 pt-4">
        <Button
          type="submit"
          variant="primary"
          size="lg"
          disabled={isLoading}
          className="min-w-[150px]"
        >
          <Typography variant="text" tag="span" className="font-medium">
            {isLoading
              ? t("common.loading")
              : mode === "create"
                ? t("events.form.submitCreate")
                : t("events.form.submitUpdate")}
          </Typography>
        </Button>
      </div>
    </form>
  );
};
