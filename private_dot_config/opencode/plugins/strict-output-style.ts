const reminder = `Strict output-style reminder (conversational replies only, never inside files):
Use exactly these five sections, with these headings and emojis, in this order, and no others:
1. 🔍 Findings/Context
2. ⚠️ Limits/Risks
3. ✅ Done
4. ✋ Recommendations
5. ➡️ Next
Findings holds only facts about the current state, never a proposed change. Recommendations holds every proposal in full (a plan or set of edits awaiting approval, even if it is the only option; which option you favor and why; optional follow-ups); Next is required actions and decisions, such as the approval request for that plan. A required decision you also have an opinion on stays as one Next item with the pick stated inline.
Omit a section only when it is empty. Do not rename, merge, or invent sections. Keep the whole block together at the end of the reply. Escape hatch: a one or two line answer may drop all five headings together, never a subset.`

type ContextEvent = {
  system: Array<{ type: "text"; text: string }>
}

type PluginContext = {
  session: {
    hook(name: "context", handler: (event: ContextEvent) => void): Promise<unknown>
  }
}

export default {
  id: "strict-output-style",

  async setup(ctx: PluginContext) {
    await ctx.session.hook("context", event => {
      event.system.push({ type: "text", text: reminder })
    })
  },
}
