declare const process: {
  env: { NONO_CAP_FILE?: string }
}

type PermissionEvent = {
  effect: "allow" | "ask" | "deny"
}

type PluginContext = {
  permission: {
    hook(name: "evaluate", handler: (event: PermissionEvent) => void): Promise<unknown>
  }
}

export default {
  id: "nono-auto-approve",

  async setup(ctx: PluginContext) {
    // This marker belongs to the server process, including ACP's private server.
    if (!process.env.NONO_CAP_FILE) return

    await ctx.permission.hook("evaluate", event => {
      // Explicit configured denials are final and never reach this hook.
      if (event.effect === "ask") event.effect = "allow"
    })
  },
}
