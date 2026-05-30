<script lang="ts">
    import "../layout.css";
    import AppSidebar from "$lib/components/AppSidebar.svelte";
    import SidebarProvider from "$lib/components/ui/sidebar/sidebar-provider.svelte";
    import SidebarTrigger from "$lib/components/ui/sidebar/sidebar-trigger.svelte";

    let { children } = $props();
</script>

<svelte:head>
    <title>Tilder | Home</title>
</svelte:head>

<SidebarProvider>
    <AppSidebar />
    <div>
        <SidebarTrigger />
        {@render children()}
    </div>
</SidebarProvider>
