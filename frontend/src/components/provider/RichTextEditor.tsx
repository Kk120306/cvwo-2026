import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import { Box, Button } from '@mui/material'
import { useEffect } from 'react'

interface RichTextEditorProps {
    content?: string
    onChange: (content: string) => void
}

const RichTextEditor = ({ content = '', onChange }: RichTextEditorProps) => {
    const editor = useEditor({
        extensions: [
            StarterKit,
            Link.configure({
                openOnClick: false,
                HTMLAttributes: {
                    class: 'text-blue-500 underline',
                },
            }),
        ],
        content,
        onUpdate: ({ editor }) => {
            onChange(editor.getHTML())
        },
    })

    useEffect(() => {
        if (editor && content !== editor.getHTML()) {
            editor.commands.setContent(content)
        }
    }, [content, editor])

    if (!editor) return null

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