import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Link, useNavigate } from "react-router";
import { Mail } from "lucide-react";

import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Field, FieldLabel, FieldError } from "~/components/ui/field";
import { toast } from "~/hooks/use-toast";
import ENV_CONFIG from "~/constants/env-config";
import { goBetterAuthClient } from "~/lib/gba-client";

const requestEmailChangeSchema = z.object({
  email: z.string().email("Please enter a valid email address"),
});

type RequestEmailChangeFormData = z.infer<typeof requestEmailChangeSchema>;

export default function RequestEmailChangePage() {
  const navigate = useNavigate();

  const form = useForm<RequestEmailChangeFormData>({
    resolver: zodResolver(requestEmailChangeSchema),
    defaultValues: {
      email: "",
    },
  });

  const onSubmit = async (data: RequestEmailChangeFormData) => {
    try {
      await goBetterAuthClient.emailPassword.requestEmailChange({
        email: data.email,
        callbackUrl: "http://localhost:3000/dashboard",
      });

      toast({
        title: "Success",
        description: "Email change requested successfully.",
      });

      navigate("/auth/sign-in");
    } catch (error: any) {
      toast({
        title: "Request email change failed",
        description: error?.message || "An unknown error occurred",
      });
    }
  };

  return (
    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
      <Controller
        name="email"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor={field.name}>New Email Address</FieldLabel>
            <div className="relative">
              <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                {...field}
                id={field.name}
                type="email"
                placeholder="newemail@example.com"
                className="pl-10"
                aria-invalid={fieldState.invalid}
              />
            </div>
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />

      <Button type="submit" className="w-full mt-6">
        Request email change
      </Button>
    </form>
  );
}
