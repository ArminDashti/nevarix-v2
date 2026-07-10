import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function ServerList() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Server List</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-muted-foreground">No servers connected.</p>
      </CardContent>
    </Card>
  );
}
