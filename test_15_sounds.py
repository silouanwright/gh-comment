#!/usr/bin/env python3
import subprocess
import time

sounds = [
    # Modern/Minimal
    ('click_soft_tap.wav', 'Very subtle tap - perfect for frequent notifications'),
    ('pop_drip.wav', 'Unique water drop - distinctive but not intrusive'),
    ('glass_ping.wav', 'Clean glass sound - modern and crisp'),
    ('click_ting_glass.wav', 'Glass click - sharp but pleasant'),

    # Pleasant/Musical
    ('music_marimba_note.wav', 'Single marimba note - warm and friendly'),
    ('chime_lite_ding_mid.wav', 'Gentle chime - soft and welcoming'),
    ('music_kalimba_on.wav', 'Kalimba sound - unique and pleasant'),
    ('chord_nice.wav', 'Pleasant chord - positive feeling'),

    # Tech/Digital
    ('digi_ping_up.wav', 'Digital ping up - futuristic'),
    ('beep_digi_note.wav', 'Digital beep - classic tech sound'),
    ('digi_blip_up.wav', 'Quick blip - minimal tech'),

    # Classic/Traditional
    ('bell_ding_hi.wav', 'Classic bell - traditional notification'),
    ('chime_done.wav', 'Completion chime - satisfying'),
    ('chime_clickbell_confirm.wav', 'Confirmation bell - clear feedback'),

    # Unique/Ambient
    ('pad_soft_on.wav', 'Ambient pad - very subtle and modern'),
]

print("\n=== Testing 15 Recommended Notification Sounds ===\n")
print("Each sound will play automatically with a 2 second gap\n")

for i, (sound, description) in enumerate(sounds, 1):
    print(f"{i}/15: {sound}")
    print(f"      {description}")

    sound_path = f"/Users/silouan.wright/Downloads/dev_tones/tones-wav/{sound}"
    subprocess.run(['afplay', sound_path])

    if i < len(sounds):
        time.sleep(2)

    print()

print("\nAll sounds played! These 15 sounds offer a good variety for different notification types.")
