"use client";

import { useAction } from 'next-safe-action/hooks';
import { useRouter } from 'next/navigation';

import { toast } from 'sonner';

import { signOutAction } from '@/app/actions';
import { Button } from '@/components/ui/button';

export function SignOutButton() {
  const router = useRouter();
  const { isPending, executeAsync } = useAction(signOutAction);

  const handleSignOut = async () => {
    try {
      const data = await executeAsync();
      if (data.serverError) {
        throw new Error(data.serverError);
      }
      if (data.validationErrors) {
        throw new Error(
          Object.entries(data.validationErrors)
            .map(([_, v]) => v)
            .join(", "),
        );
      }
      router.push("/auth/sign-in");
    } catch (error: any) {
      console.error("Error during sign in:", error);
      toast.error(error.message);
    }
  };

  return (
    <Button
      type="button"
      variant="destructive"
      className="w-full"
      disabled={isPending}
      onClick={handleSignOut}
    >
      {isPending ? "Signing out..." : "Sign Out"}
    </Button>
  );
}
