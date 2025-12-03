"use client";

import { useAction } from 'next-safe-action/hooks';
import { useRouter } from 'next/navigation';

import { signOutAction } from '@/app/actions/auth.actions';
import { Button } from '@/components/ui/button';

export function SignOutButton() {
  const router = useRouter();
  const { isPending, executeAsync } = useAction(signOutAction);

  const handleSignOut = async () => {
    await executeAsync();
    router.push("/auth/sign-in");
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
