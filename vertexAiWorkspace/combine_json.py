import json
import argparse
import os

def convert_vertexai_messages_to_md(input_file, output_file):
    """
    Converts a Vertex AI Studio conversation export (JSON file with a "messages"
    array) into a single, combined Markdown file.

    This script is designed for exports structured like:
    {
      "context": "...",
      "messages": [
        {"author": "user", "content": "Hello"},
        {"author": "model", "content": "Hi there!"}
      ]
    }

    Args:
        input_file (str): The path to the input JSON file.
        output_file (str): The path to the output Markdown file.
    """
    try:
        with open(input_file, 'r', encoding='utf-8') as infile, \
             open(output_file, 'w', encoding='utf-8') as outfile:

            outfile.write("# Vertex AI Conversation Log\n\n")

            try:
                data = json.load(infile)

                if 'messages' not in data or not isinstance(data['messages'], list):
                    print(f"Error: Input file '{input_file}' does not contain a 'messages' list.")
                    outfile.write("## Error\n\nInput file does not contain a valid 'messages' list.\n")
                    return

                messages = data['messages']
                if not messages:
                    outfile.write("*No messages found in the file.*\n")
                    print(f"Successfully created '{output_file}', but no messages were found.")
                    return

                for i, message in enumerate(messages):
                    # {
                    #   "author": "user",
                    #   "content": {
                    #     "role": "user",
                    #     "parts": [
                    #       {
                    #         "text": "BadgerDb quad store with versioning and git like metadata, what would data model look like?"
                    #       }
                    #     ]
                    #   }
                    # Determine the author/role of the message
                    if 'author' in message:
                        author = message.get('author', 'unknown').capitalize()
                    elif 'role' in message:
                        author = message.get('role', 'unknown').capitalize()
                    else:
                        author = f"Message {i+1} (Unknown Author)"
                   
                    # if author != "Bot":
                    #     continue
                    print(author)
                    # Determine the content of the message
                    if 'content' in message:
                        content = message.get('content', '')
                    elif 'text' in message:
                        content = message.get('text', '')
                    else:
                        content = "*No text content found in this message.*"

                    # {'role': 'user', 'parts': [{'text': 'BadgerDb quad store with versioning and git like metadata, what would data model look like?'}]}
                    parts = content['parts']

                    # Write the formatted message to the Markdown file
                    outfile.write(f"## {author}\n\n")
                    # outfile.write(f"{content}\n\n")
                    for part in parts:
                        # print(part['text'])
                        # Add a separator for clarity
                        outfile.write(f"{part['text']}\n")

                print(f"Successfully converted {len(messages)} messages from '{input_file}' to '{output_file}'")

            except json.JSONDecodeError:
                print(f"Error: Could not parse '{input_file}'. It may not be a valid JSON file.")
                outfile.write("## Error\n\nCould not parse the input file. Please ensure it is a valid JSON file.\n")

    except FileNotFoundError:
        print(f"Error: The input file '{input_file}' was not found.")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")

if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        description='Combine messages from a Vertex AI Studio JSON export into a single Markdown file.',
        formatter_class=argparse.RawTextHelpFormatter
    )
    parser.add_argument(
        'input_file',
        help='Path to the Vertex AI Studio export JSON file (e.g., "chat-bison-export.json").'
    )
    parser.add_argument(
        '-o', '--output',
        dest='output_file',
        help='Path for the output Markdown file (e.g., "conversation.md").\nIf not provided, it defaults to the input file name with a .md extension.'
    )

    args = parser.parse_args()

    if not args.output_file:
        base_name = os.path.splitext(args.input_file)[0]
        args.output_file = f"{base_name}.md"

    convert_vertexai_messages_to_md(args.input_file, args.output_file)
    # python3 combine_json.py BadgerDb\ Quad\ Store\ CLI\ Design.json