import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function Settings() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Settings</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-muted-foreground">Application preferences and SMTP configuration.</p>
      </CardContent>
    </Card>
  );
}
