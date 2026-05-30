<script>
    import { page } from "$app/state";
    import { Button } from "$lib/components/ui/button";
    import { MoveLeft } from "@lucide/svelte";
</script>

<svelte:head>
    <title>Tilder | {page.status}</title>
</svelte:head>

<div class="flex flex-col justify-center items-center min-h-screen">
    <h1 class="text-5xl font-bold">{page.status}</h1>
    <span class="text-3xl">{page.error?.message}</span>
    <Button
        class="cursor-pointer"
        variant="link"
        onclick={() => history.back()}
    >
        <MoveLeft />
        Go back
    </Button>
</div>
