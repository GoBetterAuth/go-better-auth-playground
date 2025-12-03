import { Mail } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function EmailVerificationPage() {
  return (
    <div className="h-full w-full p-4 grid place-items-center">
      <Card className="w-full max-w-md mx-auto mt-10">
        <CardHeader className="text-center">
          <Mail className="mx-auto h-12 w-12 text-muted-foreground" />
          <CardTitle>Check Your Email</CardTitle>
        </CardHeader>
        <CardContent className="text-center">
          <p className="text-sm text-muted-foreground mb-4">
            We&#39;ve sent a verification link to your email. Click the link to
            verify your account.
          </p>
          <Button variant="outline">Resend Email</Button>
        </CardContent>
      </Card>
    </div>
  );
}
