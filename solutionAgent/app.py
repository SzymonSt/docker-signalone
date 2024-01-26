from agent import ChatAgent
import gradio as gr

# Initialize the ChatAgent
agent = ChatAgent()

# Define the Gradio interface
iface = gr.Interface(
    fn=agent.run,
    inputs=gr.Textbox(lines=3, label="Error Logs"),
    outputs=gr.Textbox(label="Solution")
)

# Launch the interface
iface.launch()
