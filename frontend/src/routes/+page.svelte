<script lang="ts">
  import { todoClient } from "$lib/api";
  import { onMount } from "svelte";

  const loadTodos = async () => {
    let res;

    res = await todoClient.createTodo({
      text: "Hello",
    });
    console.log(res);
    const id = res.todo?.id;

    res = await todoClient.updateTodo({
      id,
      completed: true,
    });
    console.log(res);

    res = await todoClient.getTodos({});
    console.log(res);

    res = await todoClient.deleteTodo({
      id,
    });
    console.log(res);

    res = await todoClient.getTodos({});
    console.log(res);
  };

  onMount(async () => {
    await loadTodos();
  });
</script>

<h1>Welcome to SvelteKit</h1>
<p>
  Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the
  documentation
</p>
