<script>
    import { input } from '$lib/input.ts';
    import Idea from '$lib/idea.svelte';
</script>

<Idea
    title={input[0].title}
    desc={input[0].desc}
    children={input[0].children}
/>
