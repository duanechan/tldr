<script>
    import { House } from "@lucide/svelte";
    import {
        Sidebar,
        SidebarContent,
        SidebarFooter,
        SidebarHeader,
        SidebarMenu,
        SidebarMenuButton,
        SidebarMenuItem,
    } from "./ui/sidebar";
</script>

<Sidebar collapsible="icon">
    <SidebarHeader></SidebarHeader>
    <SidebarContent>
        <SidebarMenu>
            <SidebarMenuItem>
                <SidebarMenuButton class="m-auto">
                    <House />
                    Home
                </SidebarMenuButton>
            </SidebarMenuItem>
        </SidebarMenu>
    </SidebarContent>
    <SidebarFooter></SidebarFooter>
</Sidebar>
