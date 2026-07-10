import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function HubList() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Hub List</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-muted-foreground">No hubs connected.</p>
      </CardContent>
    </Card>
  );
}
