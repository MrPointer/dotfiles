const REMINDER =
  "Consider using `semble search` instead -- see AGENTS.md. Grep/Glob are only for exhaustive literal matches or confirming an exact string you already know."

const TARGET_TOOLS = new Set(["grep", "glob", "Grep", "Glob"])

export const SembleReminderPlugin = async () => {
  return {
    "tool.definition": async (input: { toolID: string }, output: { description: string }) => {
      if (!TARGET_TOOLS.has(input.toolID)) return
      if (output.description.includes(REMINDER)) return

      output.description = `${output.description}\n\nReminder: ${REMINDER}`
    },

    "tool.execute.after": async (
      input: { tool: string },
      output: { output: string },
    ) => {
      if (!TARGET_TOOLS.has(input.tool)) return
      if (output.output.includes(REMINDER)) return

      output.output = `${output.output}\n\n<system-reminder>\n${REMINDER}\n</system-reminder>`
    },
  }
}
