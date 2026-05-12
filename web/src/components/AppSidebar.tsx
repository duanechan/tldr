import { HouseIcon, SignOutIcon } from "@phosphor-icons/react";
import { Link } from "react-router-dom";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "./ui/sidebar";

export function AppSidebar() {
  function handleLogout() {
    localStorage.removeItem("ACCESS_TOKEN");
  }

  return (
    <Sidebar collapsible="icon">
      <SidebarContent className="p-2">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton className="rounded-lg" asChild>
              <Link to="/">
                <HouseIcon />
                <span>Home</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenuButton className="rounded-lg" asChild>
          <Link to="/login" onClick={handleLogout}>
            <SignOutIcon />
            <span>Logout</span>
          </Link>
        </SidebarMenuButton>
      </SidebarFooter>
    </Sidebar>
  );
}
