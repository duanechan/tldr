<script lang="ts">
    import "../layout.css";

    let { children } = $props();
</script>

<div class="flex flex-col justify-center items-center min-h-screen">
    {@render children()}
</div>
