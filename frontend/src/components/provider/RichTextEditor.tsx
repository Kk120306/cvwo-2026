import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import { Box, Button } from '@mui/material'
import { useEffect } from 'react'

// Props for content and how the change in content is handled 
interface RichTextEditorProps {
    content?: string
    onChange: (content: string) => void
}

// Rich text content editor componenet
const RichTextEditor = ({ content = '', onChange }: RichTextEditorProps) => {
    // https://tiptap.dev/docs/editor/getting-started/configure
    const editor = useEditor({
        extensions: [StarterKit, Link],
        content,
        onUpdate: ({ editor }) => {
            onChange(editor.getHTML()) // get HTML to save to backend
        },
    })

    // If content prop changes from data passed down and not thorugh editor content, we update the new content in the editor 
    useEffect(() => {
        if (editor && content !== editor.getHTML()) {
            editor.commands.setContent(content)
        }
    }, [content, editor])

    if (!editor) return null

    // https://tiptap.dev/docs/examples/basics/formatting Refer to documentation to see how its being done
    return (
        <Box>
            <Box mb={2} display="flex" gap={1}>
                <Button onClick={() => editor.chain().focus().toggleBold().run()}>Bold</Button>
                <Button onClick={() => editor.chain().focus().toggleItalic().run()}>Italic</Button>
                <Button onClick={() => editor.chain().focus().toggleBulletList().run()}>Bullets</Button>
                <Button onClick={() => editor.chain().focus().toggleOrderedList().run()}>Numbered</Button>
                <Button
                    onClick={() => {
                        let url = prompt('Enter URL')
                        if (url) {
                            // add http protocol , without it it will trail current site path 
                            if (!/^https?:\/\//i.test(url)) {
                                url = 'https://' + url
                            }
                            editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
                        }
                    }}
                >
                    Link
                </Button>
            </Box>
            <EditorContent editor={editor} style={{ border: '1px solid #ccc', padding: '10px', minHeight: '150px' }} />
        </Box>
    )
}

export default RichTextEditor
