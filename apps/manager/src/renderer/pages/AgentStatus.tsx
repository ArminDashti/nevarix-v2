import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function AgentStatus() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Agent Status</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-muted-foreground">No agents registered.</p>
      </CardContent>
    </Card>
  );
}
