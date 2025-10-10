'use client';

import { useState } from 'react';

import Composer from './Composer';
import Message from './Message';
import QuickChips from './QuickChips';

export default function ChatShell() {
  const [messages, setMessages] = useState([
    { role: 'system', content: 'Ask anything about your meetings...' },
  ]);

  return (
    <div className="mx-auto flex h-dvh max-w-screen-sm flex-col gap-3 p-4">
      <header className="text-2xl font-semibold">Ask Anything</header>
      <div className="flex-1 space-y-4 overflow-y-auto">
        {messages.map((m, i) => (
          <Message key={i} role={m.role} content={m.content} />
        ))}
      </div>
      <QuickChips
        onPick={(t) => setMessages((m) => [...m, { role: 'user', content: t }])}
      />
      <Composer
        onSend={(t) =>
          setMessages((m) => [
            ...m,
            { role: 'user', content: t },
            { role: 'assistant', content: '(streaming stub)' },
          ])
        }
      />
    </div>
  );
}
